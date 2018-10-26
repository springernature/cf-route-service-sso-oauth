package broker

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/types-cf"
)

func Catalog(w http.ResponseWriter, r *http.Request) {
	catalog := cf.Catalog{
		Services: []*cf.Service{
			{
				Name:        "sso-oauth",
				ID:          "dc081967-b7df-4bd4-8065-38415bb1e992",
				Description: "Will require users to sign in to access your app. User can simply sign in with single sign on.",
				Requires:    []string{"route_forwarding"},
				Bindable:    true,
				Plans: []*cf.Plan{
					{
						ID:          "abff1ea2-ad5e-44d4-9f00-036fcd5d16eb",
						Name:        "Google",
						Description: "SSO is implemented using Google as identity provider. Users need to sign in with a valid corporate email address.",
						Free:        true,
					},
				},
			},
		},
	}
	json, err := json.Marshal(catalog)
	if err != nil {
		fmt.Println("Um, how did we fail to marshal this catalog:")
		//fmt.Printf("%# v\n", pretty.Formatter(catalog))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
