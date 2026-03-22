package server

import (
	"context"
	"strings"

	pb "github.com/ricocynthia/brew/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BrewServer implements the gRPC BrewService
type BrewServer struct {
	pb.UnimplementedBrewServiceServer
}

// --- Seed Data ---

var remedies = []*pb.Remedy{
	{
		Id:          1,
		Name:        "Sick Prevention Tea",
		Type:        "tea",
		Description: "A warming, immune-boosting tea to drink at the first sign of sickness or as a daily preventative during cold and flu season.",
		Origin:      "Originally from doula Erica",
		Ingredients: []*pb.Ingredient{
			{Name: "Red onion", Amount: "1/4 onion", Notes: "Sliced"},
			{Name: "Ginger", Amount: "1 inch", Notes: "Fresh, sliced or grated"},
			{Name: "Cinnamon", Amount: "1 stick"},
			{Name: "Citrus peel", Amount: "1 strip", Notes: "Orange or lemon peel"},
			{Name: "Anise", Amount: "2-3 pods"},
			{Name: "Cloves", Amount: "3-4 whole cloves"},
		},
		Properties:  []string{"immune support", "antiviral", "anti-inflammatory", "warming", "cold and flu relief"},
		Preparation: "Combine all ingredients in a pot with 3-4 cups of water. Bring to a boil then reduce to a simmer for 15-20 minutes. Strain and drink warm. Add honey to taste.",
		Notes:       "Best consumed at the first sign of illness or daily during winter and flu season.",
	},
	{
		Id:          2,
		Name:        "Tea Para La Tos",
		Type:        "tea",
		Description: "A soothing cough tea made with mullein flowers, eucalyptus, and thyme. Para la tos means 'for the cough' in Spanish.",
		Ingredients: []*pb.Ingredient{
			{Name: "Mullein flowers", Amount: "1 tbsp", Notes: "Dried"},
			{Name: "Eucalyptus", Amount: "1 tsp", Notes: "Dried leaves"},
			{Name: "Thyme", Amount: "1 tsp", Notes: "Fresh or dried"},
			{Name: "Honey", Amount: "1-2 tsp", Notes: "Raw honey preferred"},
			{Name: "Lemon", Amount: "1/2 lemon", Notes: "Freshly squeezed"},
		},
		Properties:  []string{"cough relief", "respiratory support", "anti-inflammatory", "antiseptic", "soothing"},
		Preparation: "Steep mullein flowers, eucalyptus, and thyme in hot water for 10-15 minutes. Strain well. Add honey and fresh lemon juice. Drink warm and slowly.",
		Notes:       "Mullein is a powerful lung herb. Strain carefully as the fine hairs can irritate the throat.",
	},
	{
		Id:          3,
		Name:        "Chamomile y Canela",
		Type:        "tea",
		Description: "A simple, comforting blend of chamomile and cinnamon to help relax and wind down at the end of the day.",
		Ingredients: []*pb.Ingredient{
			{Name: "Chamomile flowers", Amount: "1-2 tsp", Notes: "Dried"},
			{Name: "Cinnamon", Amount: "1 stick", Notes: "Canela preferred for a softer, sweeter flavor"},
		},
		Properties:  []string{"calming", "sleep support", "relaxation", "digestive support", "warming"},
		Preparation: "Steep chamomile flowers and cinnamon stick in hot water for 5-10 minutes. Strain and drink warm. Add honey to taste.",
		Notes:       "Best enjoyed 30 minutes before bed. Canela (Mexican cinnamon) has a softer, sweeter flavor than regular cinnamon.",
	},
	{
		Id:          4,
		Name:        "Detox Bath",
		Type:        "bath",
		Description: "A deeply cleansing and detoxifying bath soak recommended during winter and flu season to support the body's natural detox process.",
		Ingredients: []*pb.Ingredient{
			{Name: "Dead sea salt", Amount: "1-2 cups"},
			{Name: "Baking soda", Amount: "1/2 cup"},
			{Name: "Apple cider vinegar", Amount: "1/2 cup", Notes: "Raw, unfiltered"},
			{Name: "Eucalyptus oil", Amount: "10-15 drops", Notes: "Essential oil"},
		},
		Properties:  []string{"detoxifying", "cleansing", "immune support", "respiratory support", "skin health"},
		Preparation: "Draw a warm bath. Add dead sea salt and baking soda and stir to dissolve. Add apple cider vinegar and eucalyptus essential oil. Soak for 20-30 minutes.",
		Notes:       "Recommended daily during winter and flu season. Drink plenty of water before and after.",
	},
}

// --- gRPC Handlers ---

// GetRemedies returns all remedies, optionally filtered by type or property
func (s *BrewServer) GetRemedies(ctx context.Context, req *pb.GetRemediesRequest) (*pb.GetRemediesResponse, error) {
	typeFilter := strings.ToLower(req.Type)
	propertyFilter := strings.ToLower(req.Property)

	var result []*pb.Remedy
	for _, remedy := range remedies {
		matchType := typeFilter == "" || strings.ToLower(remedy.Type) == typeFilter
		matchProperty := propertyFilter == ""
		if !matchProperty {
			for _, p := range remedy.Properties {
				if strings.Contains(strings.ToLower(p), propertyFilter) {
					matchProperty = true
					break
				}
			}
		}
		if matchType && matchProperty {
			result = append(result, remedy)
		}
	}

	return &pb.GetRemediesResponse{
		Count:    int32(len(result)),
		Remedies: result,
	}, nil
}

// GetRemedy returns a single remedy by ID
func (s *BrewServer) GetRemedy(ctx context.Context, req *pb.GetRemedyRequest) (*pb.Remedy, error) {
	for _, remedy := range remedies {
		if remedy.Id == req.Id {
			return remedy, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "remedy with id %d not found", req.Id)
}

// GetIngredients returns all unique ingredients across all remedies
func (s *BrewServer) GetIngredients(ctx context.Context, req *pb.GetIngredientsRequest) (*pb.GetIngredientsResponse, error) {
	seen := map[string]bool{}
	var ingredients []string
	for _, remedy := range remedies {
		for _, ing := range remedy.Ingredients {
			name := strings.ToLower(ing.Name)
			if !seen[name] {
				seen[name] = true
				ingredients = append(ingredients, ing.Name)
			}
		}
	}
	return &pb.GetIngredientsResponse{
		Count:       int32(len(ingredients)),
		Ingredients: ingredients,
	}, nil
}