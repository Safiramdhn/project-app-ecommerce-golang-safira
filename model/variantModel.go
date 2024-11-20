package model

type Variant struct {
	ID            int             `json:"id,omitempty"`
	ProductID     string          `json:"product_id,omitempty"`
	AttributeName string          `json:"atribute_name,omitempty"`
	VariantOption []VariantOption `json:"variant_options,"`
	Detail        `json:"-"`
}

type VariantOption struct {
	ID              int     `json:"id,omitempty"`
	VariantID       int     `json:"variant_id,omitempty"`
	OptionValue     string  `json:"option_value,omitempty"`
	AdditionalPrice float64 `json:"additional_price,omitempty"`
	Stock           int     `json:"stock,omitempty"`
	Detail          `json:"-"`
}

type VariantDTO struct {
	ProductID        string `json:"product_id,omitempty"`
	AttributeName    string `json:"atribute_name"`
	VariantOptionDTO `json:"variant_option"`
}

type VariantOptionDTO struct {
	VariantID       int     `json:"variant_id"`
	OptionValue     string  `json:"option_value"`
	AdditionalPrice float64 `json:"additional_price"`
	Stock           int     `json:"stock"`
}
