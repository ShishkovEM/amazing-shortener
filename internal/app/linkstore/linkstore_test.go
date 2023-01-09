package linkstore

import "testing"

const (
	testedLongURL = "http://ya.ru/"
)

func TestCreateAndGetLink(t *testing.T) {
	// Создаём и сохраняем запись об одной ссылке
	ls := NewLinkStore()
	short := ls.CreateLink(testedLongURL)

	// Получаем ссылку по short. Ничего не получаем по другому аргументу
	link, err := ls.GetLink(short)
	if err != nil {
		t.Fatal(err)
	}

	if link.Short != short {
		t.Errorf("got link.Short=%s, expected short=%s", link.Short, short)
	}

	if link.Original != testedLongURL {
		t.Errorf("got link.Original=%s, expected testedLongURL=%s", link.Original, testedLongURL)
	}

	_, err = ls.GetLink(short + "aaa")
	if err == nil {
		t.Fatal("got nil, want error")
	}
}
