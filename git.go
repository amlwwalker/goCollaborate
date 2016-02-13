package main 

import (
	"time"
	"log"
	"io/ioutil"
	git "github.com/libgit2/git2go"
)

type Repo struct {
	repo *git.Repository
	treeId *git.Oid
	branch *git.Branch
	location string
	identity string
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

var userRepo Repo
func useGitFunctions() {
		userRepo.openRepository("/Users/alex/go/src/github.com/repo/")
		// userRepo.createBranch();
		// userRepo.stageChanges()
		// userRepo.commitChanges("commited at " + t.Format("20060102150405") + "by: " + c.id)
		// userRepo.attemptMerge() //try and merge it in
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

func (r *Repo) checkoutBranch() {
	head, err := r.repo.Head()
	check(err)
	headCommit, err := r.repo.LookupCommit(head.Target())
	check(err)
	
}
//create branch for user
func (r *Repo) createBranch() {
	t := time.Now()
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
	branch, err = r.repo.CreateBranch(t.Format("20060102150405") + "_" + r.identity, headCommit, false)
	if err != nil {
		panic(err)
	}
	r.branch = branch
}
//stage changes
func (r *Repo) stageChanges(filename string) {
	log.Println("staging changes")

	//add a file to the staging area:
	idx, err := r.repo.Index()
	if err != nil {
		panic(err)
	}
	err = idx.AddByPath(filename)
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

func(r *Repo) attemptMerge() {
	branchName, _ := r.branch.Name() 
	remoteRef, err := r.repo.References.Lookup("refs/heads/"+branchName)
	if err != nil {
		panic(err)
	}
	mergeRemoteHead, err := r.repo.AnnotatedCommitFromRef(remoteRef)
	if err != nil {
		panic(err)
	}
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = mergeRemoteHead
    if err = r.repo.Merge(mergeHeads, nil, nil); err != nil {
        panic(err)
    }
}


