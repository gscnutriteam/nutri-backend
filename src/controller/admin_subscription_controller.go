package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/utils"
	"app/src/validation"

	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminSubscriptionController struct {
	SubscriptionService service.SubscriptionService
}

// Helper function to format currency
func formatCurrency(amount int) string {
	return fmt.Sprintf("Rp %d", amount)
}

func NewAdminSubscriptionController(
	subscriptionService service.SubscriptionService,
) *AdminSubscriptionController {
	return &AdminSubscriptionController{
		SubscriptionService: subscriptionService,
	}
}

// @Tags         Admin
// @Summary      Get all user subscriptions
// @Description  Returns a list of all user subscriptions with pagination
// @Produce      json
// @Security     BearerAuth
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of subscriptions"    default(10)
// @Param        status   query     string  false   "Filter by status (active, expired, pending)"
// @Router       /admin/subscriptions [get]
// @Success      200  {object}  response.SuccessWithPaginateSubscriptions
// @Failure      403  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) GetAllUserSubscriptions(ctx *fiber.Ctx) error {
	query := &validation.SubscriptionQuery{
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
		Status: ctx.Query("status", ""),
	}

	subscriptions, totalResults, err := c.SubscriptionService.GetAllUserSubscriptions(ctx, query)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithPaginateSubscriptions{
		Status:       "success",
		Message:      "User subscriptions retrieved successfully",
		Results:      subscriptions,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalResults/int64(query.Limit) + 1,
		TotalResults: totalResults,
	})
}

// @Tags         Admin
// @Summary      Get user subscription details
// @Description  Returns details of a specific user subscription
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id   path  string  true  "Subscription ID"
// @Router       /admin/subscriptions/{subscription_id} [get]
// @Success      200  {object}  response.SuccessWithSubscription
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) GetUserSubscriptionDetails(ctx *fiber.Ctx) error {
	id := ctx.Params("subscription_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Subscription ID is required")
	}

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid subscription ID format")
	}

	subscription, err := c.SubscriptionService.GetUserSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithSubscription{
		Status:  "success",
		Message: "User subscription details retrieved successfully",
		Data:    *subscription,
	})
}

// @Tags         Admin
// @Summary      Get all subscription plans
// @Description  Returns a list of all subscription plans with their users
// @Produce      json
// @Security     BearerAuth
// @Param        with_users   query  boolean  false  "Include users for each plan"
// @Router       /admin/subscription-plans [get]
// @Success      200  {object}  example.AdminSubscriptionPlansResponse
// @Failure      403  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) GetAllSubscriptionPlans(ctx *fiber.Ctx) error {
	withUsers := ctx.QueryBool("with_users", false)
	modelPlans, err := c.SubscriptionService.GetAllSubscriptionPlansWithUsers(ctx, withUsers)
	if err != nil {
		return err
	}

	// Convert model.SubscriptionPlanWithUsers to response.SubscriptionPlanWithUsers
	var responsePlans []response.SubscriptionPlanWithUsers
	for _, plan := range modelPlans {
		responsePlans = append(responsePlans, response.SubscriptionPlanWithUsers{
			ID:             plan.ID.String(),
			Name:           plan.Name,
			Price:          plan.Price,
			PriceFormatted: plan.PriceFormatted,
			Description:    plan.Description,
			AIscanLimit:    plan.AIscanLimit,
			ValidityDays:   plan.ValidityDays,
			Features:       plan.Features,
			IsActive:       plan.IsActive,
			Users:          plan.Users,
			UserCount:      plan.UserCount,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithSubscriptionPlans{
		Status:  "success",
		Message: "All subscription plans retrieved successfully",
		Data:    responsePlans,
	})
}

// @Tags         Admin
// @Summary      Update user subscription
// @Description  Updates a user subscription (plan, status, etc.)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id       path  string  true  "Subscription ID"
// @Param        request  body  validation.UpdateSubscription  true  "Update subscription data"
// @Router       /admin/subscriptions/{subscription_id} [patch]
// @Success      200  {object}  response.SuccessWithSubscription
// @Failure      400  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) UpdateUserSubscription(ctx *fiber.Ctx) error {
	id := ctx.Params("subscription_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Subscription ID is required")
	}

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid subscription ID format")
	}

	req := new(validation.UpdateSubscription)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	subscription, err := c.SubscriptionService.UpdateUserSubscription(ctx, subscriptionID, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithSubscription{
		Status:  "success",
		Message: "User subscription updated successfully",
		Data:    *subscription,
	})
}

// @Tags         Admin
// @Summary      Delete user subscription
// @Description  Deletes a user subscription
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id   path  string  true  "Subscription ID"
// @Router       /admin/subscriptions/{subscription_id} [delete]
// @Success      200  {object}  response.Common
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) DeleteUserSubscription(ctx *fiber.Ctx) error {
	id := ctx.Params("subscription_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Subscription ID is required")
	}

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid subscription ID format")
	}

	if err := c.SubscriptionService.DeleteUserSubscription(ctx, subscriptionID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Status:  "success",
		Message: "User subscription deleted successfully",
	})
}

// @Tags         Admin
// @Summary      Get transaction logs
// @Description  Returns transaction logs for a specific user subscription
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id   path  string  true  "Subscription ID"
// @Router       /admin/subscriptions/{subscription_id}/transactions [get]
// @Success      200  {object}  example.TransactionsResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) GetTransactionLogs(ctx *fiber.Ctx) error {
	id := ctx.Params("subscription_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Subscription ID is required")
	}

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid subscription ID format")
	}

	transactions, err := c.SubscriptionService.GetTransactionsBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithTransactions{
		Status:  "success",
		Message: "Transaction logs retrieved successfully",
		Data:    transactions,
	})
}

// @Tags         Admin
// @Summary      Update payment status
// @Description  Updates the payment status of a user subscription
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id       path  string  true  "Subscription ID"
// @Param        request  body  validation.UpdatePaymentStatus  true  "Update payment status data"
// @Router       /admin/subscriptions/{subscription_id}/payment-status [patch]
// @Success      200  {object}  response.SuccessWithSubscription
// @Failure      400  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) UpdatePaymentStatus(ctx *fiber.Ctx) error {
	id := ctx.Params("subscription_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Subscription ID is required")
	}

	subscriptionID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid subscription ID format")
	}

	req := new(validation.UpdatePaymentStatus)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	subscription, err := c.SubscriptionService.UpdatePaymentStatus(ctx, subscriptionID, req.Status)
	if err != nil {
		return err
	}

	// Log payment status update activity
	admin := ctx.Locals("user")
	var adminID string
	if admin != nil {
		if u, ok := admin.(*model.User); ok && u != nil {
			adminID = u.ID.String()
		}
	}

	// Generate request ID if not available
	requestID := ctx.Locals("requestID")
	var requestIDStr string
	if requestID == nil {
		requestIDStr = fmt.Sprintf("req-%s", uuid.New().String()[:8])
	} else {
		requestIDStr = requestID.(string)
	}

	utils.LogUserActivity(utils.ActivityData{
		UserID:     adminID,
		Action:     "update_payment_status",
		Resource:   "subscription",
		ResourceID: subscriptionID.String(),
		Details: map[string]interface{}{
			"status": req.Status,
		},
		RequestID:  requestIDStr,
		IPAddress:  ctx.IP(),
		UserAgent:  ctx.Get("User-Agent"),
		StatusCode: ctx.Response().StatusCode(),
	})

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithSubscription{
		Status:  "success",
		Message: "Payment status updated successfully",
		Data:    *subscription,
	})
}

// @Tags         Admin
// @Summary      Get all transaction logs
// @Description  Returns a list of all transaction logs with pagination
// @Produce      json
// @Security     BearerAuth
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of transactions"    default(10)
// @Router       /admin/transactions [get]
// @Success      200  {object}  example.TransactionsResponse
// @Failure      403  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) GetAllTransactions(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	transactions, totalResults, err := c.SubscriptionService.GetAllTransactions(ctx, page, limit)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithTransactions{
		Status:       "success",
		Message:      "All transaction logs retrieved successfully",
		Data:         transactions,
		Page:         page,
		Limit:        limit,
		TotalPages:   int64(totalResults)/int64(limit) + 1,
		TotalResults: totalResults,
	})
}

// @Tags         Admin
// @Summary      Get transaction details
// @Description  Returns details of a specific transaction
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  string  true  "Transaction ID"
// @Router       /admin/transactions/{id} [get]
// @Success      200  {object}  example.TransactionDetailResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) GetTransactionByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Transaction ID is required")
	}

	transactionID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid transaction ID format")
	}

	transaction, err := c.SubscriptionService.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithTransaction{
		Status:  "success",
		Message: "Transaction details retrieved successfully",
		Data:    *transaction,
	})
}

// @Tags         Admin
// @Summary      Get subscription plan details
// @Description  Returns details of a specific subscription plan
// @Produce      json
// @Security     BearerAuth
// @Param        plan_id   path  string  true  "Plan ID"
// @Router       /admin/subscription-plans/{plan_id} [get]
// @Success      200  {object}  response.SuccessWithSubscriptionPlan
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) GetSubscriptionPlanByID(ctx *fiber.Ctx) error {
	id := ctx.Params("plan_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Plan ID is required")
	}

	planID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid plan ID format")
	}

	plan, err := c.SubscriptionService.GetSubscriptionPlanByID(ctx, planID)
	if err != nil {
		return err
	}

	// Parse features
	var features map[string]bool
	if err := json.Unmarshal([]byte(plan.Features), &features); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error parsing plan features")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithSubscriptionPlan{
		Status:  "success",
		Message: "Subscription plan details retrieved successfully",
		Data: response.SubscriptionPlanResponse{
			ID:             plan.ID.String(),
			Name:           plan.Name,
			Price:          plan.Price,
			PriceFormatted: formatCurrency(plan.Price),
			Description:    plan.Description,
			AIscanLimit:    plan.AIscanLimit,
			ValidityDays:   plan.ValidityDays,
			Features:       features,
			IsActive:       plan.IsActive,
		},
	})
}

// @Tags         Admin
// @Summary      Update subscription plan
// @Description  Updates a subscription plan (name, price, features, etc.)
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        plan_id   path  string  true  "Plan ID"
// @Param        request  body  validation.UpdateSubscriptionPlan  true  "Update plan data"
// @Router       /admin/subscription-plans/{plan_id} [patch]
// @Success      200  {object}  response.SuccessWithSubscriptionPlan
// @Failure      400  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (c *AdminSubscriptionController) UpdateSubscriptionPlan(ctx *fiber.Ctx) error {
	id := ctx.Params("plan_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Plan ID is required")
	}

	planID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid plan ID format")
	}

	req := new(validation.UpdateSubscriptionPlan)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	plan, err := c.SubscriptionService.UpdateSubscriptionPlan(ctx, planID, req)
	if err != nil {
		return err
	}

	// Parse features
	var features map[string]bool
	if err := json.Unmarshal([]byte(plan.Features), &features); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error parsing plan features")
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithSubscriptionPlan{
		Status:  "success",
		Message: "Subscription plan updated successfully",
		Data: response.SubscriptionPlanResponse{
			ID:             plan.ID.String(),
			Name:           plan.Name,
			Price:          plan.Price,
			PriceFormatted: formatCurrency(plan.Price),
			Description:    plan.Description,
			AIscanLimit:    plan.AIscanLimit,
			ValidityDays:   plan.ValidityDays,
			Features:       features,
			IsActive:       plan.IsActive,
		},
	})
}
