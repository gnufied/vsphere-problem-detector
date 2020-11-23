// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ImageLister helps list Images.
// All objects returned here must be treated as read-only.
type ImageLister interface {
	// List lists all Images in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Image, err error)
	// Get retrieves the Image from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Image, error)
	ImageListerExpansion
}

// imageLister implements the ImageLister interface.
type imageLister struct {
	indexer cache.Indexer
}

// NewImageLister returns a new ImageLister.
func NewImageLister(indexer cache.Indexer) ImageLister {
	return &imageLister{indexer: indexer}
}

// List lists all Images in the indexer.
func (s *imageLister) List(selector labels.Selector) (ret []*v1.Image, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Image))
	})
	return ret, err
}

// Get retrieves the Image from the index for a given name.
func (s *imageLister) Get(name string) (*v1.Image, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("image"), name)
	}
	return obj.(*v1.Image), nil
}
