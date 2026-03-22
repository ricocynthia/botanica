package server

import (
	"context"
	"strings"

	pb "github.com/ricocynthia/botanica/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BotanicaServer struct {
	pb.UnimplementedBotanicaServiceServer
}

// --- Remedy Data ---

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
		Description: "A soothing cough tea made with mullein flowers, eucalyptus, and thyme. Para la tos means for the cough in Spanish.",
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
		Notes:       "Best enjoyed 30 minutes before bed. Canela has a softer, sweeter flavor than regular cinnamon.",
	},
	{
		Id:          4,
		Name:        "Detox Bath",
		Type:        "bath",
		Description: "A deeply cleansing and detoxifying bath soak recommended during winter and flu season.",
		Ingredients: []*pb.Ingredient{
			{Name: "Dead sea salt", Amount: "1-2 cups"},
			{Name: "Baking soda", Amount: "1/2 cup"},
			{Name: "Apple cider vinegar", Amount: "1/2 cup", Notes: "Raw, unfiltered"},
			{Name: "Eucalyptus oil", Amount: "10-15 drops", Notes: "Essential oil"},
		},
		Properties:  []string{"detoxifying", "cleansing", "immune support", "respiratory support", "skin health"},
		Preparation: "Draw a warm bath. Add dead sea salt and baking soda and stir to dissolve. Add apple cider vinegar and eucalyptus essential oil. Soak for 20-30 minutes.",
		Notes:       "Drink plenty of water before and after. Do not use if you have open wounds or skin irritation.",
	},
}

// --- Forageable Data ---

var forageables = []*pb.Forageable{
	{
		Id: 1, Name: "Burdock", Category: "Plant",
		Tagline:        "A common ingredient in root beer recipes",
		Properties:     []string{"Antibacterial", "Anti-cancer", "Anti-fungal", "Antioxidant", "Liver support", "Immune support", "Hormone balance"},
		Habitat:        "Roadsides, woodland edges, disturbed soils",
		Season:         "Roots in fall, leaves in spring, seeds in fall/winter",
		Parts:          "Roots, leaves, seeds",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Large fuzzy triangular leaves with wavy edges. Purple thistle-like flowers in summer that develop into brown burrs in fall. Can grow 4-10 feet tall.",
		Harvesting:     "Roots: Loosen soil, grip base and gently pull. Leaves: Cut healthy leaves closest to base. Seeds: Harvest when burrs form Velcro-like hooks.",
		Storage:        "Roots: Scrub clean, cut to size, dry thoroughly, store in airtight container away from sunlight. Leaves and seeds: Dry thoroughly and store in airtight container.",
		Warnings:       "Avoid in large amounts during pregnancy.",
		FunFact:        "Burdock is a common ingredient in traditional root beer recipes.",
	},
	{
		Id: 2, Name: "Catnip", Category: "Plant",
		Tagline:        "Not just for cats — a powerful calming herb",
		Properties:     []string{"Anti-inflammatory", "Sleep aid", "Digestive support", "Fever reducer", "Anxiety relief", "Headache relief"},
		Habitat:        "Gardens, roadsides, woodland edges",
		Season:         "When flowers are in full bloom",
		Parts:          "Leaves, flowers",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Grows 9 inches to 3 feet tall. Coarse fuzzy gray-green to medium green heart-shaped leaves with scalloped edges. Small clusters of lavender flowers with a minty aromatic smell.",
		Harvesting:     "Leaves: Cut stems with sharp scissors leaving about 3 inches. Flowers: Remove the flowering top with scissors.",
		Storage:        "Leaves: Rinse, dry, pinch from stem, store in airtight container. Flowers: Dry thoroughly and store away from sunlight.",
		Warnings:       "Generally safe. Use with caution during pregnancy.",
		FunFact:        "Catnip can be used in veterinary clinics and shelters to help lower cats stress levels.",
	},
	{
		Id: 3, Name: "Dandelion", Category: "Plant",
		Tagline:        "The lion's tooth — medicine hiding in plain sight",
		Properties:     []string{"Antibacterial", "Anti-cancer", "Anti-inflammatory", "Antiviral", "Diuretic", "Liver cleanser", "Blood purifier"},
		Habitat:        "Lawns, meadows, roadsides — almost everywhere",
		Season:         "Roots in fall, leaves in early spring, flowers in spring/early summer",
		Parts:          "Roots, leaves, flowers",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Grows up to a foot high with a deep tap root. Yellow flower head blooms April to June. Elongated leaves with highly jagged edges growing from the base.",
		Harvesting:     "Roots: Loosen soil and pull from base. Leaves: Cut with scissors leaving 3 inches. Flowers: Cut the flowering top with scissors.",
		Storage:        "Roots: Scrub, cut, dry thoroughly, store in airtight container. Leaves and flowers: Dry thoroughly and store away from sunlight.",
		Warnings:       "Generally very safe. Avoid large amounts if on blood thinners.",
		FunFact:        "Dandelion leaves have highly jagged edges said to resemble a lion's tooth — giving the plant its name.",
	},
	{
		Id: 4, Name: "Echinacea", Category: "Plant",
		Tagline:        "Nature's immune booster — belongs to the daisy family",
		Properties:     []string{"Anti-fungal", "Antiviral", "Cold and flu relief", "Immune support", "Wound healing", "UTI treatment"},
		Habitat:        "Prairies, open woodlands, roadsides",
		Season:         "Roots in fall (2nd year+ plant), leaves in early spring, flowers when in full bloom",
		Parts:          "Roots, leaves, flowers",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Perennial herb 6-24 inches tall with woody tap root. Rough hairy stems with purple or green tinge. Narrow lance-shaped leaves with three distinct veins. Purple flower heads bloom in summer.",
		Harvesting:     "Roots: Loosen soil and gently pull. Leaves: Cut stems above lowest healthy leaf pairs. Flowers: Remove flowering top with scissors.",
		Storage:        "Roots: Scrub, cut, dry thoroughly, store in airtight container. Leaves and flowers: Dry thoroughly and store away from sunlight.",
		Warnings:       "Avoid if allergic to daisies. Not recommended for long-term continuous use.",
		FunFact:        "Echinacea belongs to the daisy family.",
	},
	{
		Id: 5, Name: "Elderberry", Category: "Plant",
		Tagline:        "Flower tea doubles as a gentle eye wash",
		Properties:     []string{"Anti-cancer", "Anti-inflammatory", "Antiviral", "Rich in antioxidants", "Immune support", "Cold and flu relief"},
		Habitat:        "Forest edges, roadsides, moist lowlands",
		Season:         "Flowers in early summer, leaves in spring, berries in early autumn when fully ripe",
		Parts:          "Flowers, leaves, berries",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Wide woody shrub up to 12 feet tall. Clusters of small white flowers turning into drooping purple fruit.",
		Harvesting:     "Flowers: Cut stem 1 inch above bloom cluster. Leaves: Firmly pinch from stem. Berries: Cut main stem 1 inch above berry cluster.",
		Storage:        "All parts: Rinse in cool water, dry thoroughly, store in airtight container away from sunlight.",
		Warnings:       "Never eat raw berries, leaves, bark, or roots — they are toxic. Always cook berries before consuming.",
		FunFact:        "Elderberry flower tea can be used as a gentle eye wash for eye irritations.",
	},
	{
		Id: 6, Name: "Goldenrod", Category: "Plant",
		Tagline:        "Leaves can be cooked and eaten like spinach",
		Properties:     []string{"Antibacterial", "Anti-fungal", "Anti-inflammatory", "Antiseptic", "Diuretic", "Kidney support", "Cardiovascular support"},
		Habitat:        "Meadows, roadsides, open sunny areas",
		Season:         "Leaves in spring/summer, flowers in summer",
		Parts:          "Leaves, flowers",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Grows 2-5 feet tall. Flowers about 1/4 inch wide in tight lengthy clusters. Leaves climb the plant with slightly jagged edges and smooth texture.",
		Harvesting:     "Leaves: Firmly pinch healthy leaves from stem. Flowers: Cut at base of stem holding the flowering head.",
		Storage:        "Leaves: Rinse, dry and store in airtight container. Flowers: Dry thoroughly and store away from sunlight.",
		Warnings:       "May cause allergic reactions in people sensitive to ragweed.",
		FunFact:        "Goldenrod leaves can be cooked and eaten like spinach.",
	},
	{
		Id: 7, Name: "Motherwort", Category: "Plant",
		Tagline:        "Lion's Tail — a powerful herb for heart and hormones",
		Properties:     []string{"Anti-inflammatory", "Sleep aid", "PMS relief", "Menopause support", "Heart health", "Anxiety relief"},
		Habitat:        "Roadsides, disturbed areas, woodland edges",
		Season:         "Leaves in spring/early summer, flowers end of summer/fall",
		Parts:          "Leaves, flowers",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Upright bush up to 6.5 feet tall. Opposite leaves resembling maple or oak. Square hairy stems. Pale pink-lavender flowers bloom June through early September.",
		Harvesting:     "Cut the top third of the stems including leaves and flowers.",
		Storage:        "Leaves: Rinse, dry, pinch from stem, store in airtight container. Flowers: Dry thoroughly and store away from sunlight.",
		Warnings:       "Avoid during pregnancy. May interact with heart medications.",
		FunFact:        "Motherwort's scientific name Leonurus is the Greek word for Lion's Tail.",
	},
	{
		Id: 8, Name: "Mullein", Category: "Plant",
		Tagline:        "Nature's toilet paper — and a powerful lung herb",
		Properties:     []string{"Anti-cancer", "Anti-inflammatory", "Antiseptic", "Respiratory support", "Wound healing", "Earache treatment"},
		Habitat:        "Roadsides, disturbed ground, open sunny areas",
		Season:         "Leaves in spring/summer (2nd year most potent), flowers in late summer/early fall",
		Parts:          "Leaves, flowers",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Velvety soft plant with long large oval gray-green leaves up to 20 inches. Second year plant has erect flowering spike up to 8 feet tall with small yellow 5-petal flowers.",
		Harvesting:     "Leaves: Gently pull healthy leaves from stem. Flowers: Gently remove from stalks when in full bloom.",
		Storage:        "Leaves: Rinse, dry, store in airtight container. Flowers: Dry thoroughly and store away from sunlight.",
		Warnings:       "Generally safe. Fine leaf hairs may irritate throat if inhaled.",
		FunFact:        "Mullein is often referred to as nature's toilet paper due to its large, soft leaves.",
	},
	{
		Id: 9, Name: "Stinging Nettle", Category: "Plant",
		Tagline:        "Blanch to remove the sting — then it's pure medicine",
		Properties:     []string{"Antihistamine", "Anti-inflammatory", "Diuretic", "Allergy relief", "Arthritis relief", "Anemia support", "Hair growth"},
		Habitat:        "Moist woodland edges, riverbanks, disturbed soils",
		Season:         "Leaves in early spring, roots in late fall",
		Parts:          "Roots, leaves",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Grows 3-8 feet tall. Soft green oval to occasionally heart-shaped leaves 1-4 inches long. Both leaves and stems have stinging and non-stinging hairs.",
		Harvesting:     "Always wear gloves! Leaves: Cut plant at nodes. Roots: Wear gloves and long sleeves, loosen soil and gently pull.",
		Storage:        "Leaves: Rinse, pinch from stem, dry and store in airtight container. Roots: Scrub, cut, dry thoroughly, store away from sunlight.",
		Warnings:       "Always wear gloves when harvesting. Avoid large amounts during pregnancy.",
		FunFact:        "If blanched in hot water, the stinging hairs on nettles will no longer sting.",
	},
	{
		Id: 10, Name: "Yellow Dock", Category: "Plant",
		Tagline:        "Seeds can be pounded into flour",
		Properties:     []string{"Digestive support", "Skin health", "Cleansing", "High in iron", "Laxative", "Liver support", "Gallbladder support"},
		Habitat:        "Roadsides, fields, disturbed ground",
		Season:         "Leaves in spring, roots in early fall through winter",
		Parts:          "Seeds, roots, leaves",
		Uses:           "Tea, salves, tonics, tinctures, poultice, culinary",
		Identification: "Flower stalks grow from base about 3 feet high with small green flowers in clusters. Coarse curly-edged leaves up to 2 feet long.",
		Harvesting:     "Seeds: Rub clusters and let seeds fall into palm. Roots: Loosen soil and pull firmly. Leaves: Cut healthy leaves from base.",
		Storage:        "All parts: Dry thoroughly and store in airtight container in cool place away from sunlight.",
		Warnings:       "High in oxalates — use with caution if prone to kidney stones.",
		FunFact:        "Yellow dock seeds can be pounded into flour.",
	},
	{
		Id: 11, Name: "Chaga", Category: "Mushroom",
		Tagline:        "Black on the outside, golden on the inside",
		Properties:     []string{"Anti-cancer", "Anti-inflammatory", "Antioxidant", "Immune support", "Heart health", "Liver support"},
		Habitat:        "Grows on birch trees in cold northern climates",
		Season:         "Year-round",
		Parts:          "The growth (conk)",
		Uses:           "Tea, tinctures",
		Identification: "Grows on birch trees with black charcoal-like appearance and crusty or cracked surface. Inside is a beautiful yellowish-gold color. Never harvest from dead or fallen trees.",
		Harvesting:     "Use a sharp knife or ax to cut into it. Leave a portion attached to help the tree and chaga survive longer.",
		Storage:        "Chop into small pieces, air dry on a flat surface for a couple of days. Once dry, store in a jar.",
		Warnings:       "May interact with blood thinners and diabetes medication. High in oxalates. Never harvest dead chaga from fallen trees.",
		FunFact:        "The inside of fresh chaga is a beautiful yellowish-gold color, very different from its black exterior.",
	},
	{
		Id: 12, Name: "Morel", Category: "Mushroom",
		Tagline:        "A prized spring mushroom — use it fresh",
		Properties:     []string{"Anti-inflammatory", "Rich in antioxidants", "High in protein", "Rich in vitamin D", "Immune support"},
		Habitat:        "Woodlands and woody edges near sycamore, hickory, ash, and elm trees",
		Season:         "Spring to early summer",
		Parts:          "Cap and stem",
		Uses:           "Culinary",
		Identification: "Honeycomb-like texture on cap. Tan to brown color ranging from beige to almost black. Very hollow stem making it lighter than the cap.",
		Harvesting:     "Pinch or cut the stem just above the soil to leave the base in the ground.",
		Storage:        "Use as soon as possible — best fresh.",
		Warnings:       "Do not confuse with false morels. When in doubt do not eat it.",
		FunFact:        "The morel's hollow stem makes it much lighter than its cap — a unique trait among mushrooms.",
	},
	{
		Id: 13, Name: "Maitake", Category: "Mushroom",
		Tagline:        "Hen of the Woods — resembles a hen's ruffled feathers",
		Properties:     []string{"Anti-inflammatory", "Antioxidant", "Blood sugar regulation", "Immune support"},
		Habitat:        "Base of oak trees and other hardwoods",
		Season:         "Late summer through late fall",
		Parts:          "Entire mushroom",
		Uses:           "Culinary — meat substitute",
		Identification: "Grows in a cluster ranging from a few inches to several feet in diameter. Resembles ruffled feathers of a hen. Brownish-gray to tan on top.",
		Harvesting:     "Use a knife to cut off what you need from the top. Stems grow thick so cutting is easier than pulling.",
		Storage:        "Best cooked right away. Can refrigerate for up to one week. Only wash before use.",
		Warnings:       "Generally safe. Confirm identification carefully.",
		FunFact:        "Maitake means dancing mushroom in Japanese — said to be named because people would dance with joy upon finding it.",
	},
	{
		Id: 14, Name: "Lion's Mane", Category: "Mushroom",
		Tagline:        "Food for the brain — looks like a cheerleader pom-pom",
		Properties:     []string{"Cognitive support", "Memory improvement", "Digestive support", "Cardiovascular support", "Immune support"},
		Habitat:        "Dead or dying hardwood trees like oak, maple, or beech",
		Season:         "Late summer through fall",
		Parts:          "Entire mushroom",
		Uses:           "Culinary, supplements",
		Identification: "Round spherical shape resembling a cheerleader pom-pom. White or cream colored. Can range from a few inches to over a foot in diameter.",
		Harvesting:     "Wait until dry to harvest as it retains a lot of water. Cut the top with a knife leaving the base attached so it can regrow.",
		Storage:        "Best fresh but can refrigerate for up to one week. Use a brush to clean rather than water.",
		Warnings:       "Generally very safe. Rare allergic reactions reported.",
		FunFact:        "Lion's mane got its nickname because of its fuzzy furry appearance similar to a lion's mane.",
	},
	{
		Id: 15, Name: "Oyster Mushroom", Category: "Mushroom",
		Tagline:        "Fan-shaped and delicate — beat the bugs to it",
		Properties:     []string{"High in protein", "Rich in vitamin B", "Immune boosting", "Anti-inflammatory", "Lowers cholesterol"},
		Habitat:        "Trees, logs, and stumps — commonly on beech and aspen",
		Season:         "Late spring through early summer",
		Parts:          "Cap",
		Uses:           "Culinary — grilling, roasting, sauteing, stir-frying, soups, meat substitute",
		Identification: "Fan or oyster shaped up to 10 inches across. White gray or brown with slightly velvety texture.",
		Harvesting:     "Cut the top off leaving the base attached. If you see lots of bugs and holes leave it for the bugs.",
		Storage:        "Extremely delicate and spoils quickly. Use fresh. Can refrigerate for up to one week.",
		Warnings:       "Confirm identification carefully. Eat cooked not raw.",
		FunFact:        "If you see lots of bugs and holes in the mushroom it is no longer a good harvest — leave it for the bugs.",
	},
	{
		Id: 16, Name: "Cauliflower Mushroom", Category: "Mushroom",
		Tagline:        "Looks like lasagna noodles — found at the base of trees",
		Properties:     []string{"Anti-inflammatory", "Antioxidant", "Rich in potassium and magnesium", "Blood sugar support", "Immune support"},
		Habitat:        "Base of trees or tree roots",
		Season:         "Spring through late fall",
		Parts:          "Entire mushroom",
		Uses:           "Culinary — saute, meat substitute, pasta dishes, soups and stews, grilled",
		Identification: "Large irregular shape resembling lasagna noodles 2-6 inches in diameter. Cream to pale yellowish or tan with firm meaty texture.",
		Harvesting:     "Can sometimes pull out by hand keeping some of the base. If unsure use a knife.",
		Storage:        "Best eaten fresh. Can last about 3 days in the refrigerator.",
		Warnings:       "Confirm identification carefully before eating.",
		FunFact:        "Its irregular ruffled shape really does resemble a pile of lasagna noodles.",
	},
	{
		Id: 17, Name: "Chicken of the Woods", Category: "Mushroom",
		Tagline:        "Bright orange and bold — one of the easiest to identify",
		Properties:     []string{"Anti-inflammatory", "Antioxidant", "Rich in vitamin B", "Rich in potassium", "Immune support"},
		Habitat:        "Trunks of living dead or dying trees like oak and beech",
		Season:         "Late spring to late fall",
		Parts:          "Entire mushroom",
		Uses:           "Culinary — saute, fried, grilled, roasted, soups, stews, meat substitute",
		Identification: "Grows in clusters few inches to over a foot in diameter. Bright orange-yellow color with smooth leathery texture.",
		Harvesting:     "Inspect for bugs first. Cut at the base. Clean by brushing off debris then rinsing in cold water.",
		Storage:        "Best eaten fresh. Can last up to a week in the refrigerator.",
		Warnings:       "Some people experience GI upset especially in large quantities. Start with a small amount.",
		FunFact:        "Its bright orange-yellow color makes it one of the easiest wild mushrooms to identify confidently.",
	},
	{
		Id: 18, Name: "Chanterelle", Category: "Mushroom",
		Tagline:        "Golden and funnel-shaped — a prized find",
		Properties:     []string{"Anti-inflammatory", "Rich in antioxidants", "Rich in vitamin D", "Immune support", "Low in calories"},
		Habitat:        "Wooded areas particularly near oak and pine trees",
		Season:         "Summer through fall",
		Parts:          "Cap and stem",
		Uses:           "Culinary — sauteed, soups, stews, pickled, grilled, roasted, omelets",
		Identification: "Funnel-like shape with bright yellow to orange color. Cap 2-10 cm wide. Found growing in wooded areas near oak and pine.",
		Harvesting:     "Cut the stem near the base with a knife. Handle with care as they are fragile.",
		Storage:        "Best used right away. Can refrigerate for up to 3-5 days.",
		Warnings:       "Can be confused with the toxic jack-o-lantern mushroom. Confirm identification carefully.",
		FunFact:        "Chanterelles are one of the most prized wild mushrooms in the world featured in fine dining restaurants everywhere.",
	},
	{
		Id: 19, Name: "Giant Puffball", Category: "Mushroom",
		Tagline:        "Can grow over a foot wide — only harvest when pure white inside",
		Properties:     []string{"Anti-inflammatory", "Antioxidant", "Rich in fiber and protein", "Rich in vitamin C", "Immune support"},
		Habitat:        "Meadows, fields, woodland edges",
		Season:         "Late summer through fall",
		Parts:          "Entire mushroom",
		Uses:           "Culinary — sauteed, grilled, fried, meat substitute, soups",
		Identification: "Large white round or slightly oval shape with smooth leathery texture. Can range over a foot in diameter. Must be pure white and firm inside.",
		Harvesting:     "Confirm still white and firm. Cut at base with sharp knife. Slice and discard any discolored or soft sections.",
		Storage:        "Can refrigerate for a few days but best used immediately after harvest.",
		Warnings:       "Never eat if the inside is not pure white — discoloration means it is past its prime or potentially unsafe.",
		FunFact:        "A single giant puffball can contain trillions of spores — one of the most prolific spore producers in the mushroom world.",
	},
	{
		Id: 20, Name: "Turkey Tail", Category: "Mushroom",
		Tagline:        "Rainbow of colors — one of the most studied medicinal mushrooms",
		Properties:     []string{"Anti-cancer", "Anti-inflammatory", "Antioxidant", "Digestive support", "Immune support", "Respiratory health"},
		Habitat:        "Trunks and branches of dead or dying hardwood trees like oak",
		Season:         "Year-round",
		Parts:          "Entire mushroom",
		Uses:           "Tea, tincture, capsules, extract, powder",
		Identification: "Thin flexible fan-like shape with velvety texture 2-8 cm in diameter. Colors range from brown to green gray and shades of blue.",
		Harvesting:     "Look for healthy brightly colored turkey tail. Cut at the base with knife or scissors.",
		Storage:        "Refrigerate fresh for up to one week. Or dry thinly sliced and store in airtight container for several months.",
		Warnings:       "Generally very safe. May cause digestive upset in large amounts.",
		FunFact:        "Turkey tail is one of the most extensively researched medicinal mushrooms in the world.",
	},
}

// --- gRPC Handlers ---

func (s *BotanicaServer) GetRemedies(ctx context.Context, req *pb.GetRemediesRequest) (*pb.GetRemediesResponse, error) {
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
	return &pb.GetRemediesResponse{Count: int32(len(result)), Remedies: result}, nil
}

func (s *BotanicaServer) GetRemedy(ctx context.Context, req *pb.GetRemedyRequest) (*pb.Remedy, error) {
	for _, remedy := range remedies {
		if remedy.Id == req.Id {
			return remedy, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "remedy with id %d not found", req.Id)
}

func (s *BotanicaServer) GetIngredients(ctx context.Context, req *pb.GetIngredientsRequest) (*pb.GetIngredientsResponse, error) {
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
	return &pb.GetIngredientsResponse{Count: int32(len(ingredients)), Ingredients: ingredients}, nil
}

func (s *BotanicaServer) GetForageables(ctx context.Context, req *pb.GetForageablesRequest) (*pb.GetForageablesResponse, error) {
	categoryFilter := strings.ToLower(req.Category)
	propertyFilter := strings.ToLower(req.Property)
	var result []*pb.Forageable
	for _, f := range forageables {
		matchCategory := categoryFilter == "" || strings.ToLower(f.Category) == categoryFilter
		matchProperty := propertyFilter == ""
		if !matchProperty {
			for _, p := range f.Properties {
				if strings.Contains(strings.ToLower(p), propertyFilter) {
					matchProperty = true
					break
				}
			}
		}
		if matchCategory && matchProperty {
			result = append(result, f)
		}
	}
	return &pb.GetForageablesResponse{Count: int32(len(result)), Forageables: result}, nil
}

func (s *BotanicaServer) GetForageable(ctx context.Context, req *pb.GetForageableRequest) (*pb.Forageable, error) {
	for _, f := range forageables {
		if f.Id == req.Id {
			return f, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "forageable with id %d not found", req.Id)
}