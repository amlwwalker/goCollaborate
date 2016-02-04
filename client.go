package main

import (
	"log"
	// "fmt"
	"net/http"
	"time"
	git "github.com/libgit2/git2go"
	"github.com/gorilla/websocket"
	"io/ioutil"
)

/*
Every time a new client connects, we need to read the current contents of the file
*/

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type client struct {
	ws *websocket.Conn
	send chan []byte // Channel storing outcoming messages
	id string
}

type Repo struct {
	repo *git.Repository
	treeId *git.Oid
	branch *git.Branch
	location string
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

func doGitStuff() {
	signature := &git.Signature{
		Name: "David Calavera",
		Email: "david.calavera@gmail.com",
		When: time.Now(),
	}

    files, _ := ioutil.ReadDir("/Users/alex/go/src/github.com/repo/")
    for _, f := range files {
            log.Println(f.Name())
    }
	repo, err := git.OpenRepository("/Users/alex/go/src/github.com/repo/")
	log.Println(repo)
    if err != nil {
        panic(err)
    }

    //get the head:
    head, err := repo.Head()
	if err != nil {
		panic(err)
	}

	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		panic(err)
	}
	//create a branch
	var branch *git.Branch
	branch, err = repo.CreateBranch("whatisthename", headCommit, false)
	if err != nil {
		panic(err)
	}

	//add a file to the staging area:
	idx, err := repo.Index()
	if err != nil {
		panic(err)
	}

	err = idx.AddByPath("storage.txt")
	if err != nil {
		panic(err)
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		panic(err)
	}

	err = idx.Write()
	if err != nil {
		panic(err)
	}
	//commit the change:
	tree, err := repo.LookupTree(treeId)
	if err != nil {
		panic(err)
	}

	commitTarget, err := repo.LookupCommit(branch.Target())
	if err != nil {
		panic(err)
	}

	message := "What a day"
	_, err = repo.CreateCommit("refs/heads/whatisthename", signature, signature, message, tree, commitTarget)
	if err != nil {
		panic(err)
	}
}

//open repository
func (r *Repo) openRepository(loc string) {
	userRepo.location = loc
	repo, err := git.OpenRepository(userRepo.location)
	log.Println(repo)
    if err != nil {
        panic(err)
    }

	r.repo = repo
}
//create branch for user
func (r *Repo) createBranch(branchName string) {
    //get the head:
    head, err := r.repo.Head()
	if err != nil {
		panic(err)
	}

	headCommit, err := r.repo.LookupCommit(head.Target())
	if err != nil {
		panic(err)
	}
	//create a branch
	var branch *git.Branch
	branch, err = r.repo.CreateBranch(branchName, headCommit, false)
	if err != nil {
		panic(err)
	}
	r.branch = branch
}
//stage changes
func (r *Repo) stageChanges() {
	log.Println("staging changes")

	//add a file to the staging area:
	idx, err := r.repo.Index()
	if err != nil {
		panic(err)
	}
	err = idx.AddByPath("storage.txt")
	if err != nil {
		panic(err)
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		panic(err)
	}
	r.treeId = treeId
	err = idx.Write()
	if err != nil {
		panic(err)
	}
}
//commit the changes to the branch
func (r *Repo) commitChanges(message string) {
	log.Println("Commiting...")
	signature := &git.Signature{
		Name: "David Calavera",
		Email: "david.calavera@gmail.com",
		When: time.Now(),
	}
	//commit the change:
	tree, err := r.repo.LookupTree(r.treeId)
	if err != nil {
		panic(err)
	}

	commitTarget, err := r.repo.LookupCommit(r.branch.Target())
	if err != nil {
		panic(err)
	}

	branchName, _ := r.branch.Name()
	log.Println("commiting to: " + branchName)
	_, err = r.repo.CreateCommit("refs/heads/"+branchName, signature, signature, message, tree, commitTarget)
	if err != nil {
		panic(err)
	}
}
//merge the branch in to master

//currently (hard coded) commits a file to a branch
// func doGitStuff() {
// 	signature := &git.Signature{
// 		Name: "David Calavera",
// 		Email: "david.calavera@gmail.com",
// 		When: time.Now(),
// 	}

//     files, _ := ioutil.ReadDir("/files")
//     for _, f := range files {
//             log.Println(f.Name())
//     }
// }

//a new connection is made
func serveWs(w http.ResponseWriter, r *http.Request) {
	identity := r.URL.Query().Get("id")
	userRepo.createBranch("alex-"+identity);
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &client{
		send: make(chan []byte, maxMessageSize),
		ws: ws,
		id: identity,
	}
	log.Println("Created client: " + c.id)
	// read the content of the file when a user connects. They need the latest version from master
 	dat, err := ioutil.ReadFile(userRepo.location + "/storage.txt")
    check(err)
    log.Println(string(dat))
  	h.content = string(dat)

	h.register <- c

	go c.writePump()
	c.readPump()
}

func (c *client) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait));
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		t := time.Now()
		log.Println(c.id + " is reading: " + string(message))
		err = ioutil.WriteFile(userRepo.location + "/storage.txt", message, 0644)
		userRepo.stageChanges()
		userRepo.commitChanges("commited at " + t.Format("20060102150405") + "by: " + c.id)
		check(err)
		h.broadcast <- string(message)

	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		log.Println("WTF: " + c.id)
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *client) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}