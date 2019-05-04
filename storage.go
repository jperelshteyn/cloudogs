package main

import (
	"os"
	"io/ioutil"
	"encoding/json"

	"github.com/google/uuid"
)

type Kennel interface {
	GetOne(dog *Dog) bool
	Save(dog *Dog) error
	GetAll() []Dog
	Remove(dog *Dog) error
}

type FileStorage struct {
	filePath string
	data map[string]Dog
}

func (fs *FileStorage) read() error {
	f, err := os.Open(fs.filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	byteData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if len(byteData) > 0 {
		return json.Unmarshal(byteData, &fs.data)
	}
	return nil
}

func (fs *FileStorage) write() error {
	byteData, err := json.Marshal(&fs.data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fs.filePath, byteData, 0644)
}

func (fs *FileStorage) Init(filePath string) error {
	var f *os.File
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err = os.Create(filePath)
		defer f.Close()
		if err != nil {
			return err
		}
	}
	fs.filePath = filePath
	fs.data = make(map[string]Dog)
	return fs.read()
} 

func (fs *FileStorage) GetOne(dog *Dog) bool {
	val, ok := fs.data[dog.Id]
	dog = &val
	return ok
}

func (fs *FileStorage) Save(dog *Dog) error {
	if dog.Id == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		dog.Id = id.String()
	}
	fs.data[dog.Id] = *dog
	return fs.write()
}

func (fs *FileStorage) GetAll() []Dog {
	var dogs []Dog
	for _, val := range fs.data {
		dogs = append(dogs, val)
	}
	return dogs
}

func (fs *FileStorage) Remove(dog *Dog) error {
	delete(fs.data, dog.Id)
	return fs.write()
}
