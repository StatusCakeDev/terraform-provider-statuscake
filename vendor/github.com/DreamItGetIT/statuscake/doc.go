// Package statuscake implements a client for statuscake.com API.
//
//  // list all `Tests`
//  c, err := statuscake.New(statuscake.Auth{Username: username, Apikey: apikey})
//  if err != nil {
//    log.Fatal(err)
//  }
//
//  tests, err := c.Tests().All()
//  if err != nil {
//    log.Fatal(err)
//  }
//
//  v := url.Values{}
//  v.Set("tags", "test1,test2")
//  testsWithFilter, err := c.Tests().AllWithFilter(v)
//  if err != nil {
//    log.Fatal(err)
//  }
//
//  // delete a `Test`
//  err = c.Tests().Delete(TestID)
//
//  // create a test
//  t := &statuscake.Test{
//    WebsiteName: "Foo",
//    WebsiteURL:  "htto://example.com",
//    ... other required args...
//  }
//
//  if err = t.Validate(); err != nil {
//    log.Fatal(err)
//  }
//
//  t2, err := c.Tests().Update(t)
//  if err != nil {
//    log.Fatal(err)
//  }
//  fmt.Printf("New Test created with id: %d\n", t2.TestID)
//
//  // get Tests details
//  t, err := tt.Detail(id)
//  ...
package statuscake
