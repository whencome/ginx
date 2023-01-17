package reqs

type WelcomeRequest struct {
    Greet string `form:"greet" label:"问候语" binding:"required" binding:"required"`
}

type SayHiRequest struct {
    Name string `form:"name" label:"姓名" binding:"required" binding:"required"`
}

type SayHelloRequest struct {
    Name string `form:"name" label:"姓名" binding:"required" binding:"required"`
}
