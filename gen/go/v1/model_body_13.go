/*
 * Vendor API V1
 *
 * Apps documentation
 *
 * API version: 1.0.0
 * Contact: info@replicated.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Body13 struct {
	// Begining date formated by year-month-day ex. 2016-02-02.
	Begin string `json:"begin"`
	// Ending date formated by year-month-day ex. 2016-02-02.
	End string `json:"end"`
	// Can be set to Monthly, Quarterly, Annually, One Time, or Other to indicate interval that this billing happens.
	Frequency string `json:"frequency"`
	// LicenseType can be set to \"dev\", \"trial\", or \"prod\"
	LicenseType string `json:"license_type"`
	// Amount of money associated with this billing event.
	Revenue string `json:"revenue"`
}
