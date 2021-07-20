//Models/User.go
package Models
import (
 "github.com/alexnassif/tennis-bro/Config"
 "fmt"
)
//GetAllUsers Fetch all user data
func GetAllUsers(users *[]User) (err error) {
 if err = Config.DB.Find(users).Error; err != nil {
    
  return err
 }
 
 return nil
}
//CreateUser ... Insert New data
func CreateUser(user *User) (err error) {
 if err = Config.DB.Create(&user).Error; err != nil {
  return err
 }
 return nil
}
//GetUserByID ... Fetch only one user by Id
func GetUserByID(user *User, id string) (err error) {
 if err = Config.DB.Where("id = ?", id).First(&user).Error; err != nil {
  return err
 }
 return nil
}

func FindUserByUsername(username string, dbUser *User) (err error){
    if err = Config.DB.Where("user_name = ?", username).First(&dbUser).Error; err != nil {
        return err
    }
    return nil
}
//UpdateUser ... Update user
func UpdateUser(user *User, id string) (err error) {
 fmt.Println(user)
 Config.DB.Save(&user)
 return nil
}
//DeleteUser ... Delete user
func DeleteUser(user *User, id string) (err error) {
 Config.DB.Where("id = ?", id).Delete(&user)
 return nil
}

func AddOnlineClient(user *OnlineClient)(err error){
    if err = Config.DB.Select("ID", "UserName").Create(&user).Error; err != nil {
        return err
       }
       return nil
}

func GetAllOnlineUsers(users *[]OnlineClient) (err error) {
    if err = Config.DB.Find(&users).Error; err != nil {
        return err
    }
    return nil
}

func RemoveOnlineUser(user *OnlineClient) (err error) {

    Config.DB.Where("id = ?", user.GetId()).Delete(&user)
    return nil

}