package response

// APISuccessResponse represents the structure of a successful API response.
//
// Example:
// {
//   "success": true,
//   "data": {
//     "id": 123,
//     "name": "Samson",
//     "email": "
//   }
// }
type APISuccessResponse struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
}