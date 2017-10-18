# directory
Contact directory programming example

Basic example REST API to serve contact information.  Currently only a single API and a test file for testing each REST call.

Server: localhost
Port: 8081
API:
Get the full directory
  URI: "/directory"
  Method: "GET"
Get a single contact
  URI: "/contact/{id}"
  Method: "GET"
Add a contact
  URI: "/contact/{firstname}/{lastname}/{email}/{phone}"
  Method: "POST"
Delete a contact
  URI: "/contact/{id}"
  Method: "DELETE"

TODO:
1. Add CSV API.
2. Separate generic service logic and directory/contact specific logic.
3. Add config file for initializing directory/contact data.
4. Create separate packages for generic utilies.
