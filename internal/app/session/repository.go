package session

import (
	"context"
	"strconv"
	"time"

	"analytics-api/configs"
	str "analytics-api/internal/pkg/string"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	GetAllSession(userID, websiteID string, listSessionID []string, session session) ([]session, error)
	GetAllSessionID(userID, websiteID string, session session) ([]string, error)

	GetSessionIDToday(userID, websiteID string, session session) ([]string, error)
	GetSession(userID, sessionID string, session *session) error

	GetCountSession(userID, sessionID string) (int64, error)
	InsertSession(session session, event event) error

	GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*event, error)

	GetSessionTimestamp(sessionID string) (int64, error)
	InsertSessionTimestamp(sessionID string, timeStart int64) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetSession get session by session id
func (instance *repository) GetSession(userID, sessionID string, aSession *session) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.id": sessionID},
	}}
	err := sessionCollection.FindOne(context.TODO(), filter).Decode(&aSession)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *repository) GetAllSession(userID, websiteID string, listSessionID []string, aSession session) ([]session, error) {
	var listSession []session
	opt := options.FindOne()
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)

	for _, sessionID := range listSessionID {
		count, err := sessionCollection.CountDocuments(context.TODO(), bson.M{"$and": []bson.M{
			{"meta_data.id": sessionID},
			{"meta_data.website_id": websiteID},
			{"meta_data.user_id": userID},
		}})
		if err != nil {
			return nil, err
		}
		opt.SetSkip(count - 1)
		err = sessionCollection.FindOne(context.TODO(), bson.M{"$and": []bson.M{
			{"meta_data.id": sessionID},
			{"meta_data.website_id": websiteID},
			{"meta_data.user_id": userID},
		}}, opt).Decode(&aSession)
		if err != nil {
			return nil, err
		}
		listSession = append(listSession, aSession)
	}
	return listSession, nil
}

// GetAllSessionID get all id of session all time
func (instance *repository) GetAllSessionID(userID, websiteID string, aSession session) ([]string, error) {
	var listSessionID []string
	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.website_id": websiteID},
	}}

	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	cursor, err := sessionCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&aSession)
		if err != nil {
			return nil, err
		}
		listSessionID = append(listSessionID, aSession.MetaData.ID)
	}
	listSessionID = str.RemoveDuplicateValues(listSessionID)
	return listSessionID, nil
}

// GetAllSessionID get all id of session today
func (instance *repository) GetSessionIDToday(userID, websiteID string, aSession session) ([]string, error) {
	var listSessionID []string

	fromDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	toDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 24, 0, 0, 0, time.UTC)

	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.website_id": websiteID},
		{"time_report": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		}},
	}}

	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	cursor, err := sessionCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&aSession)
		if err != nil {
			return nil, err
		}
		listSessionID = append(listSessionID, aSession.MetaData.ID)
	}
	listSessionID = str.RemoveDuplicateValues(listSessionID)
	return listSessionID, nil
}

// InsertSession insert session
func (instance *repository) InsertSession(aSession session, event event) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	docs := session{
		MetaData: metaData{
			ID:        aSession.MetaData.ID,
			UserID:    aSession.MetaData.UserID,
			WebsiteID: aSession.MetaData.WebsiteID,
			Country:   aSession.MetaData.Country,
			City:      aSession.MetaData.City,
			Device:    aSession.MetaData.Device,
			OS:        aSession.MetaData.OS,
			Browser:   aSession.MetaData.Browser,
			Version:   aSession.MetaData.Version,
			CreatedAt: aSession.MetaData.CreatedAt,
		},
		Duration:   aSession.Duration,
		Event:      event,
		TimeReport: aSession.TimeReport,
	}
	_, err := sessionCollection.InsertOne(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}

// GetCountSession get count session of session id
func (instance *repository) GetCountSession(userID, sessionID string) (int64, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.id": sessionID},
	}}
	count, err := sessionCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *repository) GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*event, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)

	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.id": sessionID},
	}}
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip)).SetLimit(int64(limit))

	cur, err := sessionCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	var events []*event
	for cur.Next(context.TODO()) {
		var session session
		err := cur.Decode(&session)
		if err != nil {
			return nil, err
		}
		events = append(events, &session.Event)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	return events, nil
}

// InsertSessionTimestamp insert first timestamp by session id
func (instance *repository) InsertSessionTimestamp(sessionID string, timeStart int64) error {
	err := configs.Redis.Client.Set(sessionID, timeStart, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetSessionTimestamp get first timestamp by session id
func (instance *repository) GetSessionTimestamp(sessionID string) (int64, error) {
	timeStartStr, err := configs.Redis.Client.Get(sessionID).Result()
	if err != nil {
		return 0, err
	}
	timeStart, err := strconv.ParseInt(timeStartStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return timeStart, nil
}
