package updatehandler

import "encoding/xml"

type SdnList struct {
	XMLName           xml.Name `xml:"sdnList"`
	Text              string   `xml:",chardata"`
	Xsi               string   `xml:"xsi,attr"`
	Xmlns             string   `xml:"xmlns,attr"`
	PublshInformation struct {
		Text        string `xml:",chardata"`
		PublishDate string `xml:"Publish_Date"`
		RecordCount string `xml:"Record_Count"`
	} `xml:"publshInformation"`
	SdnEntry []struct {
		Text        string `xml:",chardata"`
		Uid         string `xml:"uid"`
		LastName    string `xml:"lastName"`
		SdnType     string `xml:"sdnType"`
		ProgramList struct {
			Text    string   `xml:",chardata"`
			Program []string `xml:"program"`
		} `xml:"programList"`
		AkaList struct {
			Text string `xml:",chardata"`
			Aka  []struct {
				Text      string `xml:",chardata"`
				Uid       string `xml:"uid"`
				Type      string `xml:"type"`
				Category  string `xml:"category"`
				LastName  string `xml:"lastName"`
				FirstName string `xml:"firstName"`
			} `xml:"aka"`
		} `xml:"akaList"`
		AddressList struct {
			Text    string `xml:",chardata"`
			Address []struct {
				Text            string `xml:",chardata"`
				Uid             string `xml:"uid"`
				City            string `xml:"city"`
				Country         string `xml:"country"`
				Address1        string `xml:"address1"`
				PostalCode      string `xml:"postalCode"`
				Address2        string `xml:"address2"`
				StateOrProvince string `xml:"stateOrProvince"`
				Address3        string `xml:"address3"`
			} `xml:"address"`
		} `xml:"addressList"`
		IdList struct {
			Text string `xml:",chardata"`
			ID   []struct {
				Text           string `xml:",chardata"`
				Uid            string `xml:"uid"`
				IdType         string `xml:"idType"`
				IdNumber       string `xml:"idNumber"`
				IdCountry      string `xml:"idCountry"`
				IssueDate      string `xml:"issueDate"`
				ExpirationDate string `xml:"expirationDate"`
			} `xml:"id"`
		} `xml:"idList"`
		FirstName       string `xml:"firstName"`
		Title           string `xml:"title"`
		DateOfBirthList struct {
			Text            string `xml:",chardata"`
			DateOfBirthItem []struct {
				Text        string `xml:",chardata"`
				Uid         string `xml:"uid"`
				DateOfBirth string `xml:"dateOfBirth"`
				MainEntry   string `xml:"mainEntry"`
			} `xml:"dateOfBirthItem"`
		} `xml:"dateOfBirthList"`
		PlaceOfBirthList struct {
			Text             string `xml:",chardata"`
			PlaceOfBirthItem []struct {
				Text         string `xml:",chardata"`
				Uid          string `xml:"uid"`
				PlaceOfBirth string `xml:"placeOfBirth"`
				MainEntry    string `xml:"mainEntry"`
			} `xml:"placeOfBirthItem"`
		} `xml:"placeOfBirthList"`
		NationalityList struct {
			Text        string `xml:",chardata"`
			Nationality []struct {
				Text      string `xml:",chardata"`
				Uid       string `xml:"uid"`
				Country   string `xml:"country"`
				MainEntry string `xml:"mainEntry"`
			} `xml:"nationality"`
		} `xml:"nationalityList"`
		Remarks    string `xml:"remarks"`
		VesselInfo struct {
			Text                   string `xml:",chardata"`
			CallSign               string `xml:"callSign"`
			VesselType             string `xml:"vesselType"`
			VesselFlag             string `xml:"vesselFlag"`
			VesselOwner            string `xml:"vesselOwner"`
			GrossRegisteredTonnage string `xml:"grossRegisteredTonnage"`
			Tonnage                string `xml:"tonnage"`
		} `xml:"vesselInfo"`
		CitizenshipList struct {
			Text        string `xml:",chardata"`
			Citizenship []struct {
				Text      string `xml:",chardata"`
				Uid       string `xml:"uid"`
				Country   string `xml:"country"`
				MainEntry string `xml:"mainEntry"`
			} `xml:"citizenship"`
		} `xml:"citizenshipList"`
	} `xml:"sdnEntry"`
}
