package model

type Unit struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	IsFood       *bool  `json:"is_food"`
}

type Category struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	IsFood bool   `json:"is_food"`
}

type Item struct {
	ID                   int        `json:"id"`
	Name                 string     `json:"name"`
	IsFood               bool       `json:"is_food"`
	UnitID               *int       `json:"unit_id"`
	PackageSize          *float64   `json:"package_size"`
	ExpirationDate       *string    `json:"expiration_date"`
	BestByDate           *string    `json:"best_by_date"`
	LowQuantityThreshold *float64   `json:"low_quantity_threshold"`
	CreatedAt            string     `json:"created_at"`
	UpdatedAt            string     `json:"updated_at"`
	Categories           []Category `json:"categories,omitempty"`
	Stock                []Stock    `json:"stock,omitempty"`
	TotalQuantity        float64    `json:"total_quantity"`
	Unit                 *Unit      `json:"unit,omitempty"`
}

type Stock struct {
	ID        int     `json:"id"`
	ItemID    int     `json:"item_id"`
	Location  string  `json:"location"`
	Quantity  float64 `json:"quantity"`
	DateAdded string  `json:"date_added"`
	UpdatedAt string  `json:"updated_at"`
}

type StockMoveRequest struct {
	ToLocation string  `json:"to_location"`
	Quantity   float64 `json:"quantity"`
}

type ItemFilter struct {
	Location string
	Category string
	Sort     string
	Search   string
}

type AuditEntry struct {
	ID        int     `json:"id"`
	ItemName  string  `json:"item_name"`
	Quantity  float64 `json:"quantity"`
	CreatedAt string  `json:"created_at"`
}

type Location struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	ParentID       *int       `json:"parent_id"`
	IsFood         bool       `json:"is_food"`
	ExcludeDefault bool       `json:"exclude_default"`
	Children       []Location `json:"children,omitempty"`
}
