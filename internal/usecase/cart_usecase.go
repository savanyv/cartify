package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
)

type CartUsecase struct {
	cartRepo model.CartRepository
	productVariantRepo model.ProductVariantRepository
}

func NewCartUsecase(cr model.CartRepository, pvr model.ProductVariantRepository) *CartUsecase {
	return &CartUsecase{
		cartRepo: cr,
		productVariantRepo: pvr,
	}
}

func (u *CartUsecase) GetCart(ctx context.Context, userID string) (*dto.CartResponse, error) {
	cart, err := u.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return nil, err
	}

	var items []dto.CartItemResponse
	var totalPrice float64

	for _, item := range cart.Items {
		subTotal := item.Price * float64(item.Quantity)
		totalPrice += subTotal

		productName := ""
		variantName := ""
		if item.ProductVariant.Product.ID != uuid.Nil {
			productName = item.ProductVariant.Product.Name
		}
		variantName = item.ProductVariant.Name

		items = append(items, dto.CartItemResponse{
			ID: item.ID.String(),
			ProductVariantID: item.ProductVariantID.String(),
			ProductName: productName,
			VariantName: variantName,
			Quantity: item.Quantity,
			Price: item.Price,
			SubTotal: subTotal,
		})
	}

	return &dto.CartResponse{
		ID: cart.ID.String(),
		Items: items,
		TotalPrice: totalPrice,
		ItemCount: len(items),
	}, nil
}

func (u *CartUsecase) AddToCart(ctx context.Context, userID string, req dto.AddToCartRequest) error {
	variant, err := u.productVariantRepo.FindByID(ctx, req.ProductVariantID)
	if err != nil {
		return nil
	}
	if variant == nil {
		return errors.New("product variant not found")
	}
	if variant.Stock < req.Quantity {
		return errors.New("insufficient stock")
	}

	cart, err := u.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	return u.cartRepo.AddItem(ctx, cart.ID.String(), req.ProductVariantID, variant.Price, req.Quantity)
}

func (u *CartUsecase) UpdateCartItem(ctx context.Context, userID string, itemID string, req dto.UpdateCartItemRequest) error {
	cart, err := u.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return err
	}

	var targetItem *model.CartItem
	for i, item := range cart.Items {
		if item.ID.String() == itemID {
			targetItem = &cart.Items[i]
			break
		}
	}
	if targetItem == nil {
		return errors.New("item not found in cart")
	}

	variant, err := u.productVariantRepo.FindByID(ctx, targetItem.ProductVariantID.String())
	if err != nil {
		return err
	}
	if variant == nil {
		return errors.New("product variant not found")
	}
	if variant.Stock < req.Quantity {
		return errors.New("insufficient stock")
	}

	return u.cartRepo.UpdateItemQuantity(ctx, itemID, req.Quantity)
}

func (u *CartUsecase) RemoveFromCart(ctx context.Context, userID string, itemID string) error {
	cart, err := u.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		if item.ID.String() == itemID {
			return u.cartRepo.RemoveItem(ctx, itemID)
		}
	}
	return errors.New("item not found in cart")
}

func (u *CartUsecase) ClearCart(ctx context.Context, userID string) error {
	cart, err := u.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return err
	}

	return u.cartRepo.ClearCart(ctx, cart.ID.String())
}