package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StubUserRepo struct {
	users        map[int64]User
	usersByEmail map[string]User
	nextID       int64
}

func NewStubUserRepo() *StubUserRepo {
	return &StubUserRepo{
		users:        make(map[int64]User),
		usersByEmail: make(map[string]User),
		nextID:       1,
	}
}

func (s *StubUserRepo) GetUser(id int64) (User, error) {
	user, ok := s.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return user, nil
}

func (s *StubUserRepo) GetUserByEmail(email string) (User, error) {
	user, ok := s.usersByEmail[email]
	if !ok {
		return User{}, ErrNotFound
	}
	return user, nil
}

func (s *StubUserRepo) CreateUser(u User) (int64, error) {
	u.ID = s.nextID
	s.users[s.nextID] = u
	s.usersByEmail[u.Email] = u
	s.nextID++
	return u.ID, nil
}

func (s *StubUserRepo) DeleteUser(id int64) error {
	user, ok := s.users[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.users, id)
	delete(s.usersByEmail, user.Email)
	return nil
}

func TestService_FindIDByEmail_Success(t *testing.T) {
	stub := NewStubUserRepo()
	user := User{Email: "alice@example.com", Name: "Alice"}
	id, _ := stub.CreateUser(user)

	service := New(stub)
	result, err := service.FindIDByEmail("alice@example.com")

	require.NoError(t, err)
	assert.Equal(t, id, result)
}

func TestService_FindIDByEmail_NotFound(t *testing.T) {
	stub := NewStubUserRepo()
	service := New(stub)

	_, err := service.FindIDByEmail("nonexistent@example.com")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_FindIDByEmail_InvalidEmail(t *testing.T) {
	stub := NewStubUserRepo()
	service := New(stub)

	_, err := service.FindIDByEmail("")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidEmail)
}

func TestService_FindByEmail_Table(t *testing.T) {
	stub := NewStubUserRepo()
	stub.CreateUser(User{Email: "alice@example.com", Name: "Alice"})
	stub.CreateUser(User{Email: "bob@example.com", Name: "Bob"})

	service := New(stub)

	cases := []struct {
		name      string
		email     string
		wantError bool
		wantName  string
	}{
		{"find alice", "alice@example.com", false, "Alice"},
		{"find bob", "bob@example.com", false, "Bob"},
		{"not found", "charlie@example.com", true, ""},
		{"empty email", "", true, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user, err := service.FindByEmail(c.email)

			if c.wantError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, c.wantName, user.Name)
			}
		})
	}
}

func TestService_GetUserByID_Success(t *testing.T) {
	stub := NewStubUserRepo()
	id, _ := stub.CreateUser(User{Email: "alice@example.com", Name: "Alice"})

	service := New(stub)
	user, err := service.GetUserByID(id)

	require.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
}

func TestService_GetUserByID_NotFound(t *testing.T) {
	stub := NewStubUserRepo()
	service := New(stub)

	_, err := service.GetUserByID(999)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestService_GetUserByID_InvalidID(t *testing.T) {
	stub := NewStubUserRepo()
	service := New(stub)

	invalidIDs := []int64{0, -1, -100}
	for _, id := range invalidIDs {
		_, err := service.GetUserByID(id)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidID)
	}
}

func TestService_CreateUser_Success(t *testing.T) {
	stub := NewStubUserRepo()
	service := New(stub)

	id, err := service.CreateUser("alice@example.com", "Alice")

	require.NoError(t, err)
	assert.Greater(t, id, int64(0))

	user, _ := stub.GetUser(id)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
}

func TestService_CreateUser_InvalidEmail(t *testing.T) {
	stub := NewStubUserRepo()
	service := New(stub)

	invalidEmails := []string{"", "invalid", "no@domain", "@nodomain.com"}
	for _, email := range invalidEmails {
		_, err := service.CreateUser(email, "Alice")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidEmail)
	}
}

func TestService_CreateUser_EmptyName(t *testing.T) {
	stub := NewStubUserRepo()
	service := New(stub)

	_, err := service.CreateUser("alice@example.com", "")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEmptyName)
}

func BenchmarkService_FindIDByEmail(b *testing.B) {
	stub := NewStubUserRepo()
	stub.CreateUser(User{Email: "alice@example.com", Name: "Alice"})

	service := New(stub)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.FindIDByEmail("alice@example.com")
	}
}

func BenchmarkService_GetUserByID(b *testing.B) {
	stub := NewStubUserRepo()
	id, _ := stub.CreateUser(User{Email: "alice@example.com", Name: "Alice"})

	service := New(stub)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUserByID(id)
	}
}
