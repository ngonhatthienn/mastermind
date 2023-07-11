package ToProto 

import(
	pb "intern2023/pb"
	user "intern2023/handler/User"
)

// func ToUserProto() *pb.User {
// 	var userProtos  
// }
func ToListUserProto(users []user.User) []*pb.User {
	var userProtos []*pb.User
	for _, user := range users {
		userProto := &pb.User{
			XId: user.ID,
			Username: user.Username,
			Email: user.Email,
			Password: user.Password,
			Role: user.Role,
		}
		userProtos = append(userProtos, userProto)
	}
	return userProtos
}

