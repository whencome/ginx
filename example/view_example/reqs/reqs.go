package reqs

type TestRequest struct {
    Name string `binding:"required" form:"name" label:"用户名称"`
}
