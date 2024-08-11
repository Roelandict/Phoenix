package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Social struct {
	Facebook string `json:"facebook,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
}

type User struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"`
	Age    int         `json:"age"`
	Social interface{} `json:"social"`
}

type FilteredData struct {
	Data []string `json:"data"`
}

func main() {
	usersDir := ".\\Users"
	filterDir := ".\\Filter"

	// Maak de directories als ze nog niet bestaan
	os.MkdirAll(usersDir, os.ModePerm)
	os.MkdirAll(filterDir, os.ModePerm)

	addUsers(usersDir)
	filterUsers(usersDir, filterDir)
}

func addUsers(usersDir string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		var user User

		fmt.Println("\nVoer gebruikersinformatie in (druk op Enter zonder naam in te voeren om te stoppen):")

		fmt.Print("Naam: ")
		user.Name, _ = reader.ReadString('\n')
		user.Name = strings.TrimSpace(user.Name)
		if user.Name == "" {
			break
		}

		fmt.Print("Type: ")
		user.Type, _ = reader.ReadString('\n')
		user.Type = strings.TrimSpace(user.Type)

		fmt.Print("Leeftijd: ")
		ageStr, _ := reader.ReadString('\n')
		ageStr = strings.TrimSpace(ageStr)
		user.Age, _ = strconv.Atoi(ageStr)

		fmt.Print("Heeft de gebruiker sociale media? (j/n): ")
		hasSocial, _ := reader.ReadString('\n')
		hasSocial = strings.TrimSpace(strings.ToLower(hasSocial))

		if hasSocial == "j" || hasSocial == "ja" {
			social := Social{}
			fmt.Print("Facebook: ")
			social.Facebook, _ = reader.ReadString('\n')
			social.Facebook = strings.TrimSpace(social.Facebook)

			fmt.Print("Twitter: ")
			social.Twitter, _ = reader.ReadString('\n')
			social.Twitter = strings.TrimSpace(social.Twitter)

			user.Social = social
		} else {
			user.Social = ""
		}

		// Genereer een veilige bestandsnaam
		safeFileName := strings.ReplaceAll(user.Name, " ", "_") + ".json"
		outputFile := filepath.Join(usersDir, safeFileName)

		// Schrijf de gebruiker naar een individueel JSON-bestand
		outputData, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			fmt.Printf("Fout bij het encoderen van JSON voor %s: %v\n", user.Name, err)
			continue
		}

		err = ioutil.WriteFile(outputFile, outputData, 0644)
		if err != nil {
			fmt.Printf("Fout bij het schrijven naar het bestand voor %s: %v\n", user.Name, err)
			continue
		}

		fmt.Printf("Gebruikersinformatie voor %s is succesvol toegevoegd aan %s\n", user.Name, outputFile)
	}

	fmt.Println("Gebruikers toevoegen voltooid.")
}

func filterUsers(usersDir, filterDir string) {
	files, err := ioutil.ReadDir(usersDir)
	if err != nil {
		fmt.Printf("Fout bij het lezen van de gebruikersdirectory: %v\n", err)
		return
	}

	var allUsers []User

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		data, err := ioutil.ReadFile(filepath.Join(usersDir, file.Name()))
		if err != nil {
			fmt.Printf("Fout bij het lezen van bestand %s: %v\n", file.Name(), err)
			continue
		}

		var user User
		err = json.Unmarshal(data, &user)
		if err != nil {
			fmt.Printf("Fout bij het decoderen van JSON voor %s: %v\n", file.Name(), err)
			continue
		}

		allUsers = append(allUsers, user)
	}

	filterAndWrite(allUsers, filterDir, "namen", func(u User) string { return u.Name })
	filterAndWrite(allUsers, filterDir, "types", func(u User) string { return u.Type })
	filterAndWrite(allUsers, filterDir, "leeftijden", func(u User) string { return strconv.Itoa(u.Age) })
	filterAndWriteSocial(allUsers, filterDir)

	fmt.Println("Alle categorieën zijn succesvol gefilterd en geschreven.")
}

func filterAndWrite(users []User, dir, category string, getField func(User) string) {
	var filtered FilteredData
	for _, user := range users {
		filtered.Data = append(filtered.Data, getField(user))
	}
	writeJSON(filtered, filepath.Join(dir, category+".json"))
}

func filterAndWriteSocial(users []User, dir string) {
	var filtered FilteredData
	for _, user := range users {
		switch s := user.Social.(type) {
		case map[string]interface{}:
			if fb, ok := s["facebook"].(string); ok && fb != "" {
				filtered.Data = append(filtered.Data, fmt.Sprintf("%s: Facebook: %s", user.Name, fb))
			}
			if tw, ok := s["twitter"].(string); ok && tw != "" {
				filtered.Data = append(filtered.Data, fmt.Sprintf("%s: Twitter: %s", user.Name, tw))
			}
		case string:
			if s != "" {
				filtered.Data = append(filtered.Data, fmt.Sprintf("%s: %s", user.Name, s))
			}
		}
	}
	writeJSON(filtered, filepath.Join(dir, "sociale_media.json"))
}

func writeJSON(data interface{}, filename string) {
	outputData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Fout bij het encoderen van JSON voor %s: %v\n", filename, err)
		return
	}

	err = ioutil.WriteFile(filename, outputData, 0644)
	if err != nil {
		fmt.Printf("Fout bij het schrijven naar het bestand %s: %v\n", filename, err)
		return
	}

	fmt.Printf("Gefilterde data is succesvol geschreven naar %s\n", filename)
}
