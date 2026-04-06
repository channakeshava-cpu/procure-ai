package services

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"

	"procure-ai/models"
)

type AgentService struct {
	vendorService *VendorService
	weights       models.AgentWeights
}

func NewAgentService(vendorService *VendorService) *AgentService {
	return &AgentService{
		vendorService: vendorService,
		weights: models.AgentWeights{
			Price:       0.35,
			Delivery:    0.25,
			Trust:       0.20,
			Reliability: 0.20,
		},
	}
}

func (s *AgentService) RecommendVendors(req models.ProcurementRequest) (*models.ProcurementRecommendationResponse, error) {
	vendors := s.vendorService.GetVendorsByCategory(req.Category)
	if len(vendors) == 0 {
		return &models.ProcurementRecommendationResponse{
			TopVendors:      []models.RankedVendor{},
			RejectedVendors: []models.RejectedVendor{},
			AppliedWeights:  s.weights,
			Summary:         fmt.Sprintf("No vendors found for category %q.", req.Category),
		}, nil
	}

	topN := req.TopN
	if topN == 0 {
		topN = 5
	}

	eligible := make([]models.Vendor, 0, len(vendors))
	rejected := make([]models.RejectedVendor, 0)
	for _, vendor := range vendors {
		reasons := rejectionReasons(vendor, req)
		if len(reasons) > 0 {
			rejected = append(rejected, models.RejectedVendor{
				Vendor:  vendor,
				Reasons: reasons,
			})
			continue
		}

		eligible = append(eligible, vendor)
	}

	if len(eligible) == 0 {
		return &models.ProcurementRecommendationResponse{
			TopVendors:      []models.RankedVendor{},
			RejectedVendors: rejected,
			AppliedWeights:  s.weights,
			Summary:         "No vendors satisfied the procurement constraints.",
		}, nil
	}

	minPrice, maxPrice := eligible[0].Price, eligible[0].Price
	minDelivery, maxDelivery := eligible[0].DeliveryDays, eligible[0].DeliveryDays
	for _, vendor := range eligible[1:] {
		if vendor.Price < minPrice {
			minPrice = vendor.Price
		}
		if vendor.Price > maxPrice {
			maxPrice = vendor.Price
		}
		if vendor.DeliveryDays < minDelivery {
			minDelivery = vendor.DeliveryDays
		}
		if vendor.DeliveryDays > maxDelivery {
			maxDelivery = vendor.DeliveryDays
		}
	}

	ranked := make([]models.RankedVendor, 0, len(eligible))
	for _, vendor := range eligible {
		priceScore := normalizeDescending(vendor.Price, minPrice, maxPrice)
		deliveryScore := normalizeDescending(float64(vendor.DeliveryDays), float64(minDelivery), float64(maxDelivery))
		trustScore := vendor.Trust / 5.0
		reliabilityScore := vendor.ReliabilityScore / 100.0

		finalScore := (priceScore * s.weights.Price) +
			(deliveryScore * s.weights.Delivery) +
			(trustScore * s.weights.Trust) +
			(reliabilityScore * s.weights.Reliability)

		estimatedTotal := round(reqEstTotal(vendor.Price, req.Quantity))
		breakdown := models.VendorScoreBreakdown{
			Price:       round(priceScore),
			Delivery:    round(deliveryScore),
			Trust:       round(trustScore),
			Reliability: round(reliabilityScore),
			Final:       round(finalScore),
		}

		ranked = append(ranked, models.RankedVendor{
			Vendor:         vendor,
			EstimatedTotal: estimatedTotal,
			ScoreBreakdown: breakdown,
			Reason:         buildReason(vendor, estimatedTotal, breakdown),
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].ScoreBreakdown.Final == ranked[j].ScoreBreakdown.Final {
			return ranked[i].EstimatedTotal < ranked[j].EstimatedTotal
		}
		return ranked[i].ScoreBreakdown.Final > ranked[j].ScoreBreakdown.Final
	})

	if topN > len(ranked) {
		topN = len(ranked)
	}
	top := ranked[:topN]
	for i := range top {
		top[i].Rank = i + 1
	}

	response := &models.ProcurementRecommendationResponse{
		RecommendedVendor: &top[0],
		TopVendors:        top,
		RejectedVendors:   rejected,
		AppliedWeights:    s.weights,
		Summary: fmt.Sprintf(
			"Ranked %d eligible vendors and shortlisted the top %d for category %q.",
			len(ranked),
			len(top),
			req.Category,
		),
	}

	return response, nil
}

func rejectionReasons(vendor models.Vendor, req models.ProcurementRequest) []string {
	reasons := make([]string, 0, 4)
	if req.Quantity > vendor.Stock {
		reasons = append(reasons, fmt.Sprintf("insufficient stock: requested %d, available %d", req.Quantity, vendor.Stock))
	}
	if req.Quantity < vendor.MinOrderQty {
		reasons = append(reasons, fmt.Sprintf("minimum order quantity is %d", vendor.MinOrderQty))
	}
	if req.MaxDeliveryDays > 0 && vendor.DeliveryDays > req.MaxDeliveryDays {
		reasons = append(reasons, fmt.Sprintf("delivery takes %d days", vendor.DeliveryDays))
	}
	if req.Budget > 0 && reqEstTotal(vendor.Price, req.Quantity) > req.Budget {
		reasons = append(reasons, fmt.Sprintf("estimated total %.2f exceeds budget %.2f", reqEstTotal(vendor.Price, req.Quantity), req.Budget))
	}
	if len(req.PreferredCities) > 0 && !containsFold(req.PreferredCities, vendor.Location) {
		reasons = append(reasons, fmt.Sprintf("location %q is outside preferred cities", vendor.Location))
	}
	return reasons
}

func buildReason(vendor models.Vendor, estimatedTotal float64, breakdown models.VendorScoreBreakdown) string {
	return fmt.Sprintf(
		"%s offers %.2f estimated total, %d-day delivery, %.1f/5 trust, and %.0f reliability with final score %.3f.",
		vendor.Name,
		estimatedTotal,
		vendor.DeliveryDays,
		vendor.Trust,
		vendor.ReliabilityScore,
		breakdown.Final,
	)
}

func normalizeDescending(value, minValue, maxValue float64) float64 {
	if minValue == maxValue {
		return 1
	}
	return (maxValue - value) / (maxValue - minValue)
}

func reqEstTotal(unitPrice float64, quantity int) float64 {
	return unitPrice * float64(quantity)
}

func round(value float64) float64 {
	return math.Round(value*1000) / 1000
}

func containsFold(values []string, target string) bool {
	return slices.ContainsFunc(values, func(value string) bool {
		return strings.EqualFold(value, target)
	})
}
