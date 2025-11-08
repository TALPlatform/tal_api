package api

import (
	"github.com/TALPlatform/tal_api/pkg/random"
)

var (
	testBucket  string = random.RandomName()
	firstImage  string = "seeds/assets/images/logo.webp"
	secondImage string = "seeds/assets/images/logo2.webp"
)

// func TestUploadFile(t *testing.T) {

// 	loginRequest := connect.NewRequest(&talv1.AuthLoginRequest{
// 		LoginCode:    "admin@devkit.com",
// 		UserPassword: "123456",
// 	})
// 	loginResp, err := testClient.AuthLogin(context.Background(), loginRequest)
// 	if err != nil {
// 		panic(err)
// 	}

// 	token := "Bearer " + loginResp.Msg.LoginInfo.AccessToken // Replace with a valid test token
// 	img, err := os.ReadFile(firstImage)
// 	if err != nil {
// 		panic(err)
// 	}
// 	request := connect.NewRequest(&talv1.FileCreateRequest{
// 		Path:       "img.webp",
// 		BucketName: testBucket,
// 		FileType:   "image/webp",
// 		Reader:     img,
// 	})

// 	request.Header().Set("Authorization", token)
// 	resp, err := testClient.FileCreate(context.Background(), request)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(resp)
// }
// func TestUploadFiles(t *testing.T) {
// 	img, err := os.ReadFile(firstImage)
// 	if err != nil {
// 		panic(err)
// 	}
// 	img2, err := os.ReadFile(secondImage)
// 	if err != nil {
// 		panic(err)
// 	}
// 	request := connect.NewRequest(&talv1.FileCreateBulkRequest{
// 		Files: []*talv1.FileCreateRequest{{
// 			Path:       "initial/img2.webp",
// 			BucketName: testBucket,
// 			FileType:   "image/webp",
// 			Reader:     img,
// 		},
// 			{
// 				Path:       "initial/img3.webp",
// 				BucketName: "images",
// 				FileType:   "image/webp",
// 				Reader:     img2,
// 			}}})
// 	resp, err := testClient.FileCreateBulk(context.Background(), request)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(resp)
// }
