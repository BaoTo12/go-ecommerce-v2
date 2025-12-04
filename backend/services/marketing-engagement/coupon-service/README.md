# Coupon Service

Discount and voucher management system.

## Purpose
Create, validate, and track discount coupons and vouchers.

## Key Features
- ✅ Multiple coupon types (percentage, fixed, free shipping, BOGO)
- ✅ Usage limits (global and per-user)
- ✅ Minimum order value validation
- ✅ Maximum discount caps
- ✅ Time-based validity
- ✅ Product/category restrictions
- ✅ Coupon usage tracking
- ✅ Auto-generation of unique codes

## Coupon Types
- **Percentage**: 10% off, 20% off
- **Fixed**: $5 off, $10 off
- **Free Shipping**: No shipping cost
- **BOGO**: Buy one get one

## API
- `CreateCoupon`: Generate new coupon
- `ValidateCoupon`: Check if valid for order
- `ApplyCoupon`: Apply discount to order
