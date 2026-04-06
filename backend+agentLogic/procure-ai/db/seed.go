package db

import (
	"procure-ai/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedVendors(database *gorm.DB) error {
	vendors := []models.Vendor{
		{Name: "Alpha Industrial Supply", Price: 94.50, Trust: 4.8, DeliveryDays: 2, Stock: 450, MinOrderQty: 20, Location: "Mumbai", PaymentTerms: "Net 15", ReliabilityScore: 96, Category: "electronics"},
		{Name: "Nova Parts Co", Price: 88.90, Trust: 4.2, DeliveryDays: 5, Stock: 900, MinOrderQty: 50, Location: "Pune", PaymentTerms: "Net 30", ReliabilityScore: 89, Category: "electronics"},
		{Name: "Rapid Source Logistics", Price: 102.25, Trust: 4.7, DeliveryDays: 1, Stock: 320, MinOrderQty: 10, Location: "Bengaluru", PaymentTerms: "Advance 30%", ReliabilityScore: 94, Category: "electronics"},
		{Name: "Eco Wholesale Hub", Price: 84.00, Trust: 3.9, DeliveryDays: 6, Stock: 1200, MinOrderQty: 100, Location: "Ahmedabad", PaymentTerms: "Net 45", ReliabilityScore: 82, Category: "electronics"},
		{Name: "Prime Vendor Network", Price: 97.10, Trust: 4.9, DeliveryDays: 3, Stock: 500, MinOrderQty: 25, Location: "Hyderabad", PaymentTerms: "Net 15", ReliabilityScore: 98, Category: "electronics"},
		{Name: "Sterling Components", Price: 91.75, Trust: 4.5, DeliveryDays: 4, Stock: 650, MinOrderQty: 30, Location: "Chennai", PaymentTerms: "Net 30", ReliabilityScore: 91, Category: "electronics"},
		{Name: "Summit Procurement", Price: 86.30, Trust: 4.1, DeliveryDays: 7, Stock: 1500, MinOrderQty: 150, Location: "Surat", PaymentTerms: "Net 45", ReliabilityScore: 85, Category: "electronics"},
		{Name: "Vertex Trade Links", Price: 99.80, Trust: 4.6, DeliveryDays: 2, Stock: 280, MinOrderQty: 20, Location: "Delhi", PaymentTerms: "Net 21", ReliabilityScore: 93, Category: "electronics"},
		{Name: "BluePeak Supplies", Price: 89.40, Trust: 4.3, DeliveryDays: 5, Stock: 720, MinOrderQty: 40, Location: "Noida", PaymentTerms: "Net 30", ReliabilityScore: 88, Category: "electronics"},
		{Name: "Atlas Commerce", Price: 95.20, Trust: 4.4, DeliveryDays: 3, Stock: 410, MinOrderQty: 15, Location: "Kolkata", PaymentTerms: "Net 20", ReliabilityScore: 90, Category: "electronics"},
		{Name: "Zenith MedTech Supply", Price: 148.00, Trust: 4.9, DeliveryDays: 2, Stock: 180, MinOrderQty: 5, Location: "Mumbai", PaymentTerms: "Net 10", ReliabilityScore: 97, Category: "medical"},
		{Name: "CareAxis Distributors", Price: 132.60, Trust: 4.5, DeliveryDays: 4, Stock: 260, MinOrderQty: 10, Location: "Delhi", PaymentTerms: "Net 20", ReliabilityScore: 92, Category: "medical"},
		{Name: "MediSure Wholesale", Price: 126.40, Trust: 4.0, DeliveryDays: 6, Stock: 540, MinOrderQty: 30, Location: "Jaipur", PaymentTerms: "Net 30", ReliabilityScore: 86, Category: "medical"},
		{Name: "PulseLine Vendors", Price: 141.75, Trust: 4.7, DeliveryDays: 3, Stock: 220, MinOrderQty: 8, Location: "Hyderabad", PaymentTerms: "Advance 20%", ReliabilityScore: 95, Category: "medical"},
		{Name: "Helix Health Source", Price: 119.90, Trust: 3.8, DeliveryDays: 8, Stock: 800, MinOrderQty: 50, Location: "Lucknow", PaymentTerms: "Net 45", ReliabilityScore: 81, Category: "medical"},
		{Name: "ForgeRaw Materials", Price: 72.20, Trust: 4.3, DeliveryDays: 5, Stock: 2000, MinOrderQty: 200, Location: "Nagpur", PaymentTerms: "Net 30", ReliabilityScore: 87, Category: "raw_materials"},
		{Name: "IronBridge Traders", Price: 69.85, Trust: 4.0, DeliveryDays: 7, Stock: 2600, MinOrderQty: 300, Location: "Jamshedpur", PaymentTerms: "Net 45", ReliabilityScore: 84, Category: "raw_materials"},
		{Name: "Titan SourceWorks", Price: 76.40, Trust: 4.8, DeliveryDays: 3, Stock: 1100, MinOrderQty: 120, Location: "Vadodara", PaymentTerms: "Net 15", ReliabilityScore: 96, Category: "raw_materials"},
		{Name: "CoreBulk Industries", Price: 71.60, Trust: 3.9, DeliveryDays: 9, Stock: 3200, MinOrderQty: 500, Location: "Raipur", PaymentTerms: "Net 60", ReliabilityScore: 80, Category: "raw_materials"},
		{Name: "Harbor Line Supply", Price: 74.95, Trust: 4.6, DeliveryDays: 4, Stock: 1450, MinOrderQty: 150, Location: "Visakhapatnam", PaymentTerms: "Net 21", ReliabilityScore: 93, Category: "raw_materials"},
	}

	return database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&vendors).Error
}
