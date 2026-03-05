package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	// "go.uber.org/zap"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/utils"
)

// AuthJWTMiddleware for JWT-based authentication
func (mw *MiddlewareManager) AuthJWTMiddleware(authUC auth.UseCase, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")
			mw.logger.Infof("AuthJWTMiddleware RequestID: %s, BearerHeader: %s", utils.GetRequestId(c), bearerHeader)

			if bearerHeader == "" || !strings.HasPrefix(bearerHeader, "Bearer ") {
				mw.logger.Errorf("AuthJWTMiddleware RequestID: %s, Error: missing or invalid Authorization header", utils.GetRequestId(c))
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError("missing or invalid Authorization header"))
			}

			tokenString := strings.TrimPrefix(bearerHeader, "Bearer ")
			if err := mw.validateJWTToken(tokenString, authUC, c, cfg); err != nil {
				mw.logger.Errorf("AuthJWTMiddleware RequestID: %s, Error: %s", utils.GetRequestId(c), err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err.Error()))
			}

			return next(c)
		}
	}
}

// AdminMiddleware checks if user is admin
func (mw *MiddlewareManager) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok || user.Role == nil || *user.Role != "admin" {
			mw.logger.Errorf("AdminMiddleware RequestID: %s, Error: user is not admin", utils.GetRequestId(c))
			return c.JSON(http.StatusForbidden, httpErrors.NewForbiddenError("permission denied"))
		}
		return next(c)
	}
}

// OwnerOrAdminMiddleware allows admin or resource owner
func (mw *MiddlewareManager) OwnerOrAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*models.User)
			if !ok {
				mw.logger.Errorf("OwnerOrAdminMiddleware RequestID: %s, Error: invalid user context", utils.GetRequestId(c))
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError("invalid user context"))
			}

			if user.Role != nil && *user.Role == "admin" {
				return next(c)
			}

			if user.Id.String() != c.Param("user_id") {
				mw.logger.Errorf("OwnerOrAdminMiddleware RequestID: %s, UserID: %s, Error: user is not owner", utils.GetRequestId(c), user.Id.String())
				return c.JSON(http.StatusForbidden, httpErrors.NewForbiddenError("permission denied"))
			}

			return next(c)
		}
	}
}

// RoleBasedAuthMiddleware checks for specific roles
func (mw *MiddlewareManager) RoleBasedAuthMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*models.User)
			if !ok {
				mw.logger.Errorf("RoleBasedAuthMiddleware RequestID: %s, Error: invalid user context", utils.GetRequestId(c))
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError("invalid user context"))
			}

			if user.Role == nil {
				mw.logger.Errorf("RoleBasedAuthMiddleware RequestID: %s, UserID: %s, Error: user has no role", utils.GetRequestId(c), user.Id.String())
				return c.JSON(http.StatusForbidden, httpErrors.NewForbiddenError("permission denied"))
			}

			for _, role := range roles {
				if role == *user.Role {
					return next(c)
				}
			}

			mw.logger.Errorf("RoleBasedAuthMiddleware RequestID: %s, UserID: %s, Error: user role %s not allowed", utils.GetRequestId(c), user.Id.String(), *user.Role)
			return c.JSON(http.StatusForbidden, httpErrors.NewForbiddenError("permission denied"))
		}
	}
}

// validateJWTToken validates JWT token and sets user in context
func (mw *MiddlewareManager) validateJWTToken(tokenString string, authUC auth.UseCase, c echo.Context, cfg *config.Config) error {
	if tokenString == "" {
		return httpErrors.NewUnauthorizedError("empty JWT token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Server.JwtSecretKey), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.NewUnauthorizedError("invalid JWT token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["id"].(string); ok && strings.TrimSpace(userID) != "" {
			userUUID, err := uuid.Parse(userID)
			if err != nil {
				return err
			}

			user, err := authUC.GetByID(c.Request().Context(), userUUID)
			if err != nil {
				return err
			}

			if !user.Active {
				return httpErrors.NewUnauthorizedError("user is inactive")
			}

			c.Set("user", user)
			ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
			c.SetRequest(c.Request().WithContext(ctx))
			return nil
		}

		if userAddress, ok := claims["user_address"].(string); ok && strings.TrimSpace(userAddress) != "" {
			// Support agency wallet tokens that only carry user_address claim.
			user := &models.User{
				Id:          uuid.Nil,
				UserAddress: &userAddress,
				Active:      true,
			}
			c.Set("user", user)
			ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
			c.SetRequest(c.Request().WithContext(ctx))
			return nil
		}

		return httpErrors.NewUnauthorizedError("invalid JWT claims")
	}

	return httpErrors.NewUnauthorizedError("invalid JWT token")
}
