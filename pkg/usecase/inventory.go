package usecase

import (
	"errors"
	helper_interface "jerseyhub/pkg/helper/interface"
	interfaces "jerseyhub/pkg/repository/interface"
	"jerseyhub/pkg/utils/models"
	"mime/multipart"
)

type inventoryUseCase struct {
	repository      interfaces.InventoryRepository
	offerRepository interfaces.OfferRepository
	helper          helper_interface.Helper
}

func NewInventoryUseCase(repo interfaces.InventoryRepository, offer interfaces.OfferRepository, h helper_interface.Helper) *inventoryUseCase {
	return &inventoryUseCase{
		repository:      repo,
		offerRepository: offer,
		helper:          h,
	}
}

func (i *inventoryUseCase) AddInventory(inventory models.AddInventories, image *multipart.FileHeader) (models.InventoryResponse, error) {

	url, err := i.helper.AddImageToS3(image)
	if err != nil {
		return models.InventoryResponse{}, err
	}

	//send the url and save it in database
	InventoryResponse, err := i.repository.AddInventory(inventory, url)
	if err != nil {
		return models.InventoryResponse{}, err
	}

	return InventoryResponse, nil

}

func (i *inventoryUseCase) UpdateInventory(pid int, stock int) (models.InventoryResponse, error) {

	result, err := i.repository.CheckInventory(pid)
	if err != nil {

		return models.InventoryResponse{}, err
	}

	if !result {
		return models.InventoryResponse{}, errors.New("there is no inventory as you mentioned")
	}

	newcat, err := i.repository.UpdateInventory(pid, stock)
	if err != nil {
		return models.InventoryResponse{}, err
	}

	return newcat, err
}

func (i *inventoryUseCase) DeleteInventory(inventoryID string) error {

	err := i.repository.DeleteInventory(inventoryID)
	if err != nil {
		return err
	}
	return nil

}

func (i *inventoryUseCase) ShowIndividualProducts(id string) (models.Inventories, error) {

	product, err := i.repository.ShowIndividualProducts(id)
	if err != nil {
		return models.Inventories{}, err
	}

	DiscountPercentage, err := i.offerRepository.FindDiscountPercentage(product.CategoryID)
	if err != nil {
		return models.Inventories{}, err
	}

	//make discounted price by calculation
	var discount float64
	if DiscountPercentage > 0 {
		discount = (product.Price * float64(DiscountPercentage)) / 100
	}

	product.DiscountedPrice = product.Price - discount

	return product, nil

}

func (i *inventoryUseCase) ListProducts(page int) ([]models.Inventories, error) {

	productDetails, err := i.repository.ListProducts(page)
	if err != nil {
		return []models.Inventories{}, err
	}

	//loop inside products and then calculate discounted price of each then return
	for j := range productDetails {
		discount_percentage, err := i.offerRepository.FindDiscountPercentage(productDetails[j].CategoryID)
		if err != nil {
			return []models.Inventories{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64

		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

	}

	return productDetails, nil

}

func (i *inventoryUseCase) SearchProducts(key string) ([]models.Inventories, error) {

	productDetails, err := i.repository.SearchProducts(key)
	if err != nil {
		return []models.Inventories{}, err
	}

	//loop inside products and then calculate discounted price of each then return
	for j := range productDetails {
		discount_percentage, err := i.offerRepository.FindDiscountPercentage(productDetails[j].CategoryID)
		if err != nil {
			return []models.Inventories{}, errors.New("there was some error in finding the discounted prices")
		}
		var discount float64

		if discount_percentage > 0 {
			discount = (productDetails[j].Price * float64(discount_percentage)) / 100
		}

		productDetails[j].DiscountedPrice = productDetails[j].Price - discount

	}

	return productDetails, nil

}
