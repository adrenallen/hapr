package users

type FirebaseMessagingDeviceService struct {
	Repository *FirebaseMessagingDeviceRepository
}

func NewFirebaseMessagingDeviceService() *FirebaseMessagingDeviceService {
	s := new(FirebaseMessagingDeviceService)
	s.Repository = &FirebaseMessagingDeviceRepository{}
	return s
}

func (s *FirebaseMessagingDeviceService) RecordDeviceToUser(deviceID string, userID int) (*FirebaseMessagingDevice, error) {
	fmd, err := s.Repository.GetFMDByDeviceID(deviceID)
	if err != nil {
		//device needs new recording
		fmd = &FirebaseMessagingDevice{
			DeviceID: deviceID,
			UserID:   userID,
		}

		//save new and return
		return s.Repository.SaveNewFMD(fmd)
	}

	//update and save
	fmd.UserID = userID
	return s.Repository.UpdateFMD(fmd)
}
