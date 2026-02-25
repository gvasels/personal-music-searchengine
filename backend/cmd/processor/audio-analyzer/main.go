// Audio Analyzer Lambda - Combines signal analysis and GenAI analysis
// Uses Bedrock for genre/mood classification and Marengo for embeddings
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type Event struct {
	TrackID string `json:"trackId"`
	UserID  string `json:"userId"`
	S3Key   string `json:"s3Key"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
}

type AnalysisResult struct {
	TrackID  string   `json:"trackId"`
	UserID   string   `json:"userId"`
	Analysis Analysis `json:"analysis"`
	Error    string   `json:"error,omitempty"`
}

type Analysis struct {
	Genre           string    `json:"genre"`
	SubGenre        string    `json:"subGenre"`
	Mood            string    `json:"mood"`
	ToneDescription string    `json:"toneDescription"`
	Sections        []Section `json:"sections"`
	Instrumentation string    `json:"instrumentation"`
	VocalPresence   string    `json:"vocalPresence"`
	EnergyProfile   string    `json:"energyProfile"`
}

type Section struct {
	Name        string `json:"name"`
	StartSec    int    `json:"startSec"`
	EndSec      int    `json:"endSec"`
	Description string `json:"description"`
}

var bedrockClient *bedrockruntime.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load AWS config: %v", err))
	}
	bedrockClient = bedrockruntime.NewFromConfig(cfg)
}

func handler(ctx context.Context, event Event) (AnalysisResult, error) {
	result := AnalysisResult{
		TrackID: event.TrackID,
		UserID:  event.UserID,
	}

	// Build prompt for genre/mood analysis
	prompt := buildAnalysisPrompt(event)

	// Call Bedrock Claude for analysis
	analysis, err := analyzeWithBedrock(ctx, prompt)
	if err != nil {
		result.Error = fmt.Sprintf("bedrock analysis failed: %v", err)
		return result, nil // Return partial result, don't fail pipeline
	}

	result.Analysis = analysis
	return result, nil
}

func buildAnalysisPrompt(event Event) string {
	return fmt.Sprintf(`Analyze this music track and provide structured metadata.

Track info:
- Title: %s
- Artist: %s
- Album: %s

Based on the artist, title, and album information, infer the likely genre, mood, and characteristics.

Provide JSON output with these fields:
{
  "genre": "primary genre (e.g., Electronic, Rock, Hip Hop, Jazz, Classical)",
  "subGenre": "sub-genre (for Electronic: Deep House, Tech House, Techno, Trance, etc.)",
  "mood": "one-word mood (e.g., energetic, melancholic, uplifting, dark, chill)",
  "toneDescription": "2-3 sentence description of the likely tone and feel",
  "sections": [
    {"name": "intro", "startSec": 0, "endSec": 30, "description": "typical intro description"},
    {"name": "main", "startSec": 30, "endSec": 180, "description": "main section description"}
  ],
  "instrumentation": "likely instruments/sounds (e.g., synthesizers, guitar, drums)",
  "vocalPresence": "none|male|female|mixed",
  "energyProfile": "description of energy arc"
}

Use the genre as the top-level category and subGenre for the specific style. Use these taxonomies:

ELECTRONIC:
- House: Deep House, Tech House, Progressive House, Electro House, Future House, Afro House, Soulful House, Acid House, Jackin House, Funky House, Tribal House, Chicago House, Garage House, Bass House, Lo-Fi House, Organic House, Melodic House, Minimal House, Microhouse, Latin House, Euro House, Beach House, Vocal House, Filter House, Piano House, French House, Italo House
- Techno: Minimal Techno, Industrial Techno, Melodic Techno, Acid Techno, Detroit Techno, Dub Techno, Hard Techno, Peak Time Techno, Hypnotic Techno, Raw Techno, Atmospheric Techno, Bleep Techno, Tribal Techno, Ambient Techno, EBM, Berlin Techno, Schranz, Modular Techno
- Trance: Progressive Trance, Uplifting Trance, Psytrance, Vocal Trance, Tech Trance, Goa Trance, Hard Trance, Balearic Trance, Dream Trance, Full-On Psytrance, Dark Psytrance, Acid Trance, Epic Trance, Nitzhonot, Suomisaundi, Zenonesque, Psychill, Forest Psytrance, Hi-Tech Psytrance, Twilight Psytrance
- Drum and Bass: Liquid DnB, Neurofunk, Jump-Up, Jungle, Darkstep, Drumstep, Atmospheric DnB, Dancefloor DnB, Minimal DnB, Ragga Jungle, Techstep, Halftime DnB, Intelligent DnB, Autonomic, Soulful DnB, Deep DnB, Rollers, Amens, Cross-Breed DnB
- Dubstep: Brostep, Riddim, Deep Dubstep, Future Bass, Melodic Dubstep, Tearout, Colour Bass, Post-Dubstep, Dungeon Sound, Purple Sound, Deathstep, Chillstep, Trapstep, Hybrid Trap
- Ambient: Dark Ambient, Ambient Dub, Drone, Space Ambient, Ambient Techno, New Age, Ambient House, Chillout, Healing, Soundscape, Organic Ambient, Cosmic Ambient
- Breakbeat: Breaks, Big Beat, Nu Skool Breaks, Breakcore, Acid Breaks, Florida Breaks, Progressive Breaks, Atmospheric Breaks, Broken Beat, Uk Breaks
- Hardcore: Happy Hardcore, Gabber, Frenchcore, Speedcore, UK Hardcore, Hardstyle, Rawstyle, Terrorcore, Uptempo, Industrial Hardcore, Mainstream Hardcore, Early Hardcore, J-Core, Makina, Freeform, Reverse Bass, Euphoric Hardstyle
- Downtempo: Trip Hop, Lounge, Chillwave, Synthwave, Retrowave, Vaporwave, Lo-Fi Beats, Electronica, Psybient, Psychill, Downtempo Bass, Organic Downtempo, Glitch Hop, Midtempo Bass
- UK Bass: UK Garage, 2-Step, Grime, Bassline, UK Funky, Speed Garage, Future Garage, Night Bass, Deep Garage, 4x4 Garage, Dubwise
- Disco & Funk: Nu Disco, Disco House, Italo Disco, Cosmic Disco, Boogie, Electro Funk, Space Disco, Euro Disco, Hi-NRG, Afro Disco, Disco Polo
- IDM: Glitch, Braindance, Experimental Electronic, Microsound, Generative, Granular, Algorithmic
- Melodic House & Techno: Indie Dance, Melodic Progressive, Afro Melodic, Progressive Melodic, Organic Progressive
- Electro: Electro Clash, Miami Bass, Electro Breaks, Electro Pop, Detroit Electro, Freestyle
- Bass Music: Wave, Experimental Bass, Deconstructed Club, Footwork, Juke, Jersey Club, Baltimore Club, Baile Funk, Leftfield Bass, Future Beats, Wonky, Trap (EDM), Festival Trap
- Garage & Speed: Garage, Breakstep, 4x4
- Ethnic Electronic: Afro Tech, Afro Electronic, Middle Eastern Electronic, Indian Electronic, Latin Electronic, Tropical Bass, Global Bass, Ethno-Electro
- Hardwave & Cyberpunk: Hardwave, Darksynth, Cyberpunk, Aggrotech, Futurepop, Terrorwave

HIP HOP & RAP:
- Hip Hop: Boom Bap, East Coast, West Coast, Southern Hip Hop, Conscious Hip Hop, Underground Hip Hop, Abstract Hip Hop, Instrumental Hip Hop, Chopped & Screwed, G-Funk, Gangsta Rap, Horrorcore, Crunk
- Trap: Atlanta Trap, Drill, UK Drill, Phonk, Melodic Trap, Dark Trap, Cloud Rap, Plugg
- R&B: Contemporary R&B, Neo Soul, Alternative R&B, Quiet Storm, New Jack Swing, PBR&B, Funk R&B

ROCK:
- Rock: Classic Rock, Hard Rock, Soft Rock, Blues Rock, Southern Rock, Heartland Rock, Arena Rock, Glam Rock, Garage Rock, Psychedelic Rock, Space Rock, Stoner Rock, Surf Rock, Krautrock
- Alternative Rock: Indie Rock, Shoegaze, Dream Pop, Post-Punk, Britpop, Madchester, Noise Rock, Math Rock, Emo, Screamo, Post-Rock, Slowcore, Sadcore
- Metal: Heavy Metal, Thrash Metal, Death Metal, Black Metal, Doom Metal, Power Metal, Progressive Metal, Symphonic Metal, Metalcore, Deathcore, Djent, Nu Metal, Sludge Metal, Post-Metal, Folk Metal, Gothic Metal
- Punk: Punk Rock, Hardcore Punk, Post-Punk, Pop Punk, Skate Punk, Street Punk, Crust Punk, D-Beat, Anarcho-Punk

POP:
- Pop: Synth Pop, Electropop, Dance Pop, Eurodance, Teen Pop, Bubblegum Pop, Art Pop, Chamber Pop, Baroque Pop, Indie Pop, Dream Pop, Hyperpop, K-Pop, J-Pop, Latin Pop, City Pop

JAZZ:
- Jazz: Bebop, Cool Jazz, Hard Bop, Free Jazz, Modal Jazz, Jazz Fusion, Smooth Jazz, Acid Jazz, Nu Jazz, Swing, Big Band, Latin Jazz, Gypsy Jazz, Afro-Cuban Jazz, Spiritual Jazz, Jazz Funk, Ethio-Jazz

CLASSICAL:
- Classical: Baroque, Classical Period, Romantic, Modern Classical, Contemporary Classical, Minimalism, Neo-Classical, Orchestral, Chamber Music, Opera, Choral, Film Score

SOUL & FUNK:
- Soul: Northern Soul, Southern Soul, Philly Soul, Psychedelic Soul, Blue-Eyed Soul, Deep Soul, Gospel
- Funk: P-Funk, Deep Funk, Funk Rock, Afrobeat, Go-Go, Boogie

REGGAE & CARIBBEAN:
- Reggae: Roots Reggae, Dub, Dancehall, Ragga, Lovers Rock, Ska, Rocksteady, Soca, Calypso, Reggaeton, Dembow

LATIN:
- Latin: Salsa, Bachata, Merengue, Cumbia, Bossa Nova, Samba, Tango, Latin Jazz, Corridos, Regional Mexican, Norteño, Banda, Urbano

COUNTRY & FOLK:
- Country: Classic Country, Outlaw Country, Country Pop, Americana, Alt-Country, Bluegrass, Honky-Tonk, Country Rock, Bro-Country
- Folk: Traditional Folk, Contemporary Folk, Indie Folk, Freak Folk, Neofolk, Celtic, World Folk

BLUES:
- Blues: Delta Blues, Chicago Blues, Electric Blues, Country Blues, Blues Rock, Jump Blues, Boogie-Woogie, Swamp Blues, Texas Blues, British Blues

WORLD & GLOBAL:
- African: Afrobeats, Amapiano, Highlife, Soukous, Kwaito, Gqom, Mbalax, Jùjú, Ethio-Jazz
- Middle Eastern: Arabic Pop, Rai, Turkish Pop, Persian Pop, Andalusian
- Asian: Bollywood, Bhangra, C-Pop, Enka, Gamelan, Qawwali
- Caribbean: Soca, Zouk, Kompa, Bouyon

EXPERIMENTAL:
- Experimental: Noise, Industrial, Musique Concrète, Sound Collage, Field Recordings, Plunderphonics, Dark Ambient, Power Electronics

Return ONLY valid JSON, no other text.`, event.Title, event.Artist, event.Album)
}

func analyzeWithBedrock(ctx context.Context, prompt string) (Analysis, error) {
	modelID := os.Getenv("BEDROCK_MODEL_ID")
	if modelID == "" {
		modelID = "global.anthropic.claude-sonnet-4-6"
	}

	// Build request body for Claude
	requestBody := map[string]interface{}{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        1024,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return Analysis{}, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := bedrockClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &modelID,
		Body:        bodyBytes,
		ContentType: stringPtr("application/json"),
	})
	if err != nil {
		return Analysis{}, fmt.Errorf("invoke model: %w", err)
	}

	// Parse response
	var response struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return Analysis{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(response.Content) == 0 {
		return Analysis{}, fmt.Errorf("empty response from model")
	}

	// Parse the JSON from the response, stripping markdown code fences if present
	text := strings.TrimSpace(response.Content[0].Text)
	if strings.HasPrefix(text, "```") {
		// Remove opening fence (```json or ```)
		if idx := strings.Index(text, "\n"); idx != -1 {
			text = text[idx+1:]
		}
		// Remove closing fence
		if idx := strings.LastIndex(text, "```"); idx != -1 {
			text = text[:idx]
		}
		text = strings.TrimSpace(text)
	}

	var analysis Analysis
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		return Analysis{}, fmt.Errorf("parse analysis JSON: %w", err)
	}

	return analysis, nil
}

func stringPtr(s string) *string {
	return &s
}

func main() {
	lambda.Start(handler)
}
