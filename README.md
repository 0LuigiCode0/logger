# logger - Красивый и информативный логгер

Пример формата 

`%{font:b}%{color:f}%{color:b} %{level} %{color:rb} %{time:2006/01/02--15:04:05} >>> WHERE \n\t%{file:1}\n MSG %{r}\"%{message}\"`

* `%{font:b}` - установка жирности текста
* `%{font:rb}` - снятие жирности текста
* `%{color:f}` - установка цвета шрифта
* `%{color:rf}` - установка цвета шрифта по умолчанию
* `%{color:b}` - установка цвета фона
* `%{color:rb}` - установка цвета фона по умолчанию
* `%{r}` - Сбрасывает все настройки
* `%{level}` - Название уровня
* `%{time:2006/01/02--15:04:05}` - Время сообщения и его формат
* `%{file:1}` - Выаодит trace с уровнем отображения пути файлов
* `%{message}` - Сообщение пользователя