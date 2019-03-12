package localFileStructure

import (
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type ImageNode struct {
	Path         string
	CategoryName string
	ModTime      time.Time
	Md5Sum       string
	ImageId      int
}

func GetImageList(fileSystem map[string]*FilesystemNode) ([]*ImageNode, error) {
	logrus.Debugln("Starting GetImageList to prepare local image metadata.")

	imageFiles := make([]*ImageNode, 0, len(fileSystem))

	finished := make(chan bool, 1)
	errChannel := make(chan error, 1)
	queue := make(chan *FilesystemNode, 100)
	results := make(chan *ImageNode, 200)
	waitGroup := sync.WaitGroup{}

	go resultCollector(results, &imageFiles)

	waitGroup.Add(1)
	go queueProducer(fileSystem, queue, &waitGroup)

	numberOfCPUs := runtime.NumCPU()
	for i := 0; i < numberOfCPUs; i++ {
		logrus.Tracef("Starting getImageNodeWorker number %d", i)
		waitGroup.Add(1)
		go getImageNodeWorker(queue, results, errChannel, &waitGroup)
	}

	go func() {
		waitGroup.Wait()

		logrus.Debugln("All workers finished processing, closing channels.")

		close(results)
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			logrus.Errorf("Error during local image processing: %s", err)
			return nil, err
		}
	}

	logrus.Infof("Found %d local images to process", len(imageFiles))

	return imageFiles, nil
}

func queueProducer(fileSystem map[string]*FilesystemNode, queue chan *FilesystemNode, waitGroup *sync.WaitGroup) {
	logrus.Debugln("Starting queueProducer to fill the queue of the files to check and calculate the checksum")
	for _, file := range fileSystem {
		if file.IsDir {
			continue
		}
		queue <- file
	}

	// after the last item is in the queue, we close it as there will be no more and we like
	// the workers to exit.
	close(queue)

	logrus.Debugln("Finished queueProducer")

	waitGroup.Done()
}

func resultCollector(results chan *ImageNode, imageFiles *[]*ImageNode) {
	logrus.Debugln("Starting image node result collector")
	for imageNode := range results {
		logrus.Debugf("Local Image prepared - %s - %s - %s", imageNode.Md5Sum, imageNode.ModTime.Format(time.RFC3339), imageNode.Path)
		*imageFiles = append(*imageFiles, imageNode)
	}
	logrus.Debugln("Finished resultCollector")

}

func getImageNodeWorker(queue chan *FilesystemNode, results chan *ImageNode, errChannel chan error, waitGroup *sync.WaitGroup) {
	logrus.Debugln("Starting image file worker to gather local image informations")
	for file := range queue {
		md5sum, err := calculateFileCheckSums(file.Path)
		if err != nil {
			errChannel <- err
			// we try the next image in the queue, as this might be just one error
			continue
		}

		imageNode := &ImageNode{
			Path:         file.Path,
			CategoryName: filepath.Dir(file.Key),
			ModTime:      file.ModTime,
			Md5Sum:       md5sum,
		}

		results <- imageNode
	}

	logrus.Debugln("Finished getImageNodeWorker")
	waitGroup.Done()
}
