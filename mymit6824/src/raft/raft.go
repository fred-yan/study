package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import (
	"labrpc"
	"time"
	"math/rand"
	"github.com/CodisLabs/codis/pkg/utils/log"
)

// import "bytes"
// import "encoding/gob"



//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

type LogEntry struct {
	Term int
	Command interface{}
}

type ServerState string

const (
	Follower  ServerState = "Follower"
	Candidate             = "Candidate"
	Leader                = "Leader"
)

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	state ServerState

	//Persistent state on all servers
	currentTerm int
	voteFor     int
	log         []LogEntry

	//Volatile state on all servers
	commitIndex int
	lastApplide int

	//Volatile state on leaders
	nextIndex  []int
	matchIndex []int

	// when receive the heartbeat
	heartBeatTime time.Time
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	var term int
	var isleader bool
	// Your code here (2A).
	term = rf.currentTerm
	isleader = rf.state == Leader
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := gob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := gob.NewDecoder(r)
	// d.Decode(&rf.xxx)
	// d.Decode(&rf.yyy)
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
}




//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	log.Info("id %d receive a request vote from candidate %d and term is %d", rf.me, args.CandidateId, args.Term)
	// Reply false if term < currentTerm
	if args.Term < rf.currentTerm {
		return
	}

	// If votedFor is null or candidateId, and candidate’s log is at
	// least as up-to-date as receiver’s log, grant vot
	if rf.voteFor == -1 || rf.voteFor == args.CandidateId {
		rfLastIndex := len(rf.log) -1
		if rf.log[rfLastIndex].Term > args.LastLogTerm {
			log.Info("id %d refused vote for it's log last entry's term %d bigger than request's log term %d",
				rf.me, rf.log[rfLastIndex].Term, args.LastLogTerm)
			return
		}

		if rf.log[rfLastIndex].Term == args.LastLogTerm {
			if rfLastIndex > args.LastLogIndex {
				log.Info("id %d refused vote for it's log last entry's index %d bigger than request's log index %d",
					rf.me, rfLastIndex, args.LastLogIndex)
				return
			}
		}
	}

	log.Info("id %d grant the request vote from candidate %d and update it's term to %d",
		rf.me, args.CandidateId, args.Term)
	rf.currentTerm = args.Term
	rf.state = Follower
	rf.voteFor = args.CandidateId

	reply.Term = args.Term
	reply.VoteGranted = true
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, vote chan int, args *RequestVoteArgs, reply *RequestVoteReply)  {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	if ok {
		vote <- server
	}
}


//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).


	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.state = Follower
	rf.currentTerm = 0
	rf.voteFor = -1
	rf.lastApplide = 0
	rf.commitIndex = 0

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())


	return rf
}


func (rf *Raft) startElectionProcess() {
	electionTimeOut := (200 + time.Duration(rand.Intn(300)))*time.Millisecond
	currentTime := <- time.After(electionTimeOut)

	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.state != Leader && currentTime.Sub(rf.heartBeatTime) > electionTimeOut {
		log.Info("election timeout %f", electionTimeOut.Seconds())
		go rf.beginElection()
	}
}


func (rf *Raft) beginElection() {
	rf.mu.Lock()
	rf.state = Candidate
	rf.currentTerm++
	rf.voteFor = rf.me

	log.Info("begin election")
	lastLogIndex := len(rf.log) - 1
	lastLogTerm := rf.log[lastLogIndex].Term
	args := RequestVoteArgs{
		Term: rf.currentTerm,
		CandidateId: rf.voteFor,
		LastLogIndex: lastLogIndex,
		LastLogTerm: lastLogTerm,
	}

	replies := make([]RequestVoteReply, len(rf.peers))
	voteChan := make(chan int, len(rf.peers))
	for i:=0; i<len(rf.peers); i++ {
		if i != rf.me {
			go rf.sendRequestVote(i, voteChan, &args, &replies[i])
		}
	}
	rf.mu.Unlock()

	votes := 1
	for i := 0; i < len(replies); i++ {
		id := <-voteChan
		reply := replies[id]
		rf.mu.Lock()
		if reply.Term > rf.currentTerm {
			log.Info("id %d change from candidater to fellower of id %d", rf.me, id)
			rf.currentTerm = reply.Term
			rf.state = Follower
			rf.voteFor = -1
		} else if reply.VoteGranted {
			votes += 1
			if votes > len(replies)/2{
				if rf.state == Candidate && rf.currentTerm == args.Term {
					log.Info("id %d start to promote to leader for term %d", rf.me, rf.currentTerm)
					go rf.promoteToLeader()
					break
				} else {
					log.Info("id %d election for term %d was interrupted", rf.me, args.Term)
				}
			}
		}
		rf.mu.Unlock()
	}

}

func (rf *Raft) promoteToLeader() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.state = Leader

	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))

	for i := range rf.peers {
		if i != rf.me {
			rf.nextIndex[i] = len(rf.log) + 1 // Should be initialized to leader's last log index + 1
			rf.matchIndex[i] = 0              // Index of highest log entry known to be replicated on server

			// Start routines for each peer which will be used to monitor and send log entries
			go rf.startLeaderPeerProcess(i)
		}
	}
}


const LeaderPeerTickInterval = 10 * time.Millisecond
func (rf *Raft) startLeaderPeerProcess(peerIndex int) {
	ticker := time.NewTicker(LeaderPeerTickInterval)

	// Initial heartbeat

	lastEntrySent := time.Now()

	for {


		select {
		case currentTime := <-ticker.C: // If traffic has been idle, we should send a heartbeat
			if currentTime.Sub(lastEntrySent) >= HeartBeatInterval {
				lastEntrySent = time.Now()

			}
		}
	}
}