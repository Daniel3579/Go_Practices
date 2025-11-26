package notes

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("note not found")

type Repo struct {
	col *mongo.Collection
}

func NewRepo(db *mongo.Database) (*Repo, error) {
	col := db.Collection("notes")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "title", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "content", Value: "text"},
			},
			Options: options.Index().SetName("text_search"),
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}

	_, err := col.Indexes().CreateMany(context.Background(), indexes)

	if err != nil {
		return nil, err
	}
	return &Repo{col: col}, nil
}

func (r *Repo) Create(ctx context.Context, title, content string, expireTime *int64) (Note, error) {
	now := time.Now()
	n := Note{
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if expireTime != nil {
		expireDuration := time.Duration(*expireTime) * time.Second
		expirationTime := now.Add(expireDuration)
		n.ExpiresAt = &expirationTime
	}

	res, err := r.col.InsertOne(ctx, n)
	if err != nil {
		return Note{}, err
	}
	n.ID = res.InsertedID.(primitive.ObjectID)
	return n, nil
}

func (r *Repo) ByID(ctx context.Context, idHex string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}
	var n Note
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&n); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return n, nil
}

func (r *Repo) List(ctx context.Context, q string, after string, limit int64, skip int64) ([]Note, error) {
	filter := bson.M{}
	if q != "" {
		filter["$text"] = bson.M{"$search": q}
	}

	if after != "" {
		afterID, err := primitive.ObjectIDFromHex(after)
		if err == nil {
			filter["_id"] = bson.M{"$lt": afterID}
		} else {
			return nil, err
		}
	}

	opts := options.
		Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "_id", Value: -1}})

	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, cur.Err()
}

func (r *Repo) Update(ctx context.Context, idHex string, title, content *string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Note{}, ErrNotFound
	}

	set := bson.M{"updatedAt": time.Now()}
	if title != nil {
		set["title"] = *title
	}
	if content != nil {
		set["content"] = *content
	}

	after := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Note
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, after).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Note{}, ErrNotFound
		}
		return Note{}, err
	}
	return updated, nil
}

func (r *Repo) Delete(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return ErrNotFound
	}
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repo) GetStats(ctx context.Context) (Stats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"notesCount":    bson.M{"$sum": 1},
			"avgContentLen": bson.M{"$avg": bson.M{"$strLenCP": "$content"}},
		}}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return Stats{}, err
	}
	defer cur.Close(ctx)

	var result Stats
	if cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			return Stats{}, err
		}
	}

	return result, cur.Err()
}
