package server

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type SessionManager struct {
	lock        sync.Mutex // protects session
	maxlifetime int64
	provider    *MemProvider
}

func InitializeSessionManager() *SessionManager {
	// 86400 => 60sec * 60min * 24h
	s_mng := &SessionManager{provider: NewMemProvider(),
		maxlifetime: 86400}
	return s_mng
}

const (
	PLAYING   = 0
	AVAILABLE = 1
)

type Player struct {
	Uname  string
	Status int
	Id     string
	Ip     string
}

func (s_mng *SessionManager) SessionStart(player *Player, response *[]Player) error {
	// lock session for data syncronization problems
	s_mng.lock.Lock()
	defer s_mng.lock.Unlock()
	if !s_mng.provider.SessionExist(player.Id) {
		s_mng.provider.SessionInit(player)
	}

	s_mng.provider.getAllUsers(response)
	return nil
}

func (s_mng *SessionManager) SelectPlayer(players *[]string, response *string) error {
	s_mng.lock.Lock()
	defer s_mng.lock.Unlock()
	if s_mng.provider.SessionExist((*players)[0]) &&
		s_mng.provider.SessionExist((*players)[1]) {
		session := s_mng.provider.SessionRead((*players)[0])
		session.player.Status = PLAYING
		fmt.Println("Select: ", session.player)
		session = s_mng.provider.SessionRead((*players)[1])
		session.player.Status = PLAYING
		fmt.Println("Select: ", session.player)
	} else {
		return errors.New("No Users with thoose names: " + (*players)[0] + " and " + (*players)[1])
	}
	return nil
}

type MemProvider struct {
	lock        sync.RWMutex // ReadWrite Mutex
	session_map map[string]*MemSession
	maxlifetime int64
}

type MemSession struct {
	sid          string       // unique session id
	player       Player       // user
	timeAccessed time.Time    // last access time
	lock         sync.RWMutex // ReadWrite Mutex
}

// implements initialization of session, it returns new session
// variable if it succeed.
func (p *MemProvider) SessionInit(player *Player) *MemSession {
	p.lock.Lock() // Lock locks for writing.
	defer p.lock.Unlock()

	newSession := &MemSession{sid: player.Id,
		player:       *player,
		timeAccessed: time.Now()}

	p.session_map[player.Id] = newSession

	return newSession
}

// returns session variable that is represented by corresponding
// session_id, it creates a new session variable and return if it does not exist.
func (p *MemProvider) SessionRead(sid string) *MemSession {
	p.lock.RLock()
	go p.SessionUpdate(sid)
	p.lock.RUnlock()
	return p.session_map[sid]
}

func (p *MemProvider) getAllUsers(response *[]Player) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	for _, v := range p.session_map {
		fmt.Println(v)
		if v.player.Status == AVAILABLE {
			*response = append(*response, v.player)
		}
	}
}

func (p *MemProvider) SessionExist(sid string) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if _, exist := p.session_map[sid]; exist {
		return true
	} else {
		return false
	}
}

// deletes session variable by corresponding session_id.
func (p *MemProvider) SessionDestroy(sid string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exist := p.session_map[sid]; exist {
		delete(p.session_map, sid)
	}
}

func (p *MemProvider) SessionUpdate(sid string) {
	p.lock.Lock() // Lock locks for writing.
	defer p.lock.Unlock()
	if session, exist := p.session_map[sid]; exist {
		session.timeAccessed = time.Now()
	}
}

func NewMemProvider() *MemProvider {
	return &MemProvider{session_map: make(map[string]*MemSession),
		maxlifetime: 86400}
}
