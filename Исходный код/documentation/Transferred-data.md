# Перед работой с API
Для работы API требуется верификация клиента. О том, как верифицироваться и получить доступ к API, откройте страницу [о верификации клиентов](http://distet.tech/VKR/server/wiki/Verification-Client)
# Передаваемые данные
Здесь описаны все передаваемые данных и их формы, облегчая разработку клиентского приложения.
## Данные авторизации
Авторизация достигается двумя способами. 
### Авторизация посредством пароля и логина
Для авторизации с помощью логина и пароля необходимо отправить по пути **/Auth** JSON следующего формата:
```
{
	"email": "look example",
	"password": "look example",
	"sessia": {
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": "look example"
		},
		"device": {
# Уникальный код системы.Для linux - /etc/machine - id.Для Windows ключ machine-id HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\SQMClient.
			"deviceid": "look example",
			"os": "look example",
			"mac": "look example",
# Сетевое имя системы.Для Windows как правило DESKTOP - xxxxxx.Для linux смотри /etc/hostname 
			"hostname": "look example"
		}
	}
}
```
Возвращаемое значение - **токен**.
###### Токен – используется для всех последующих операций, начиная от отправки сообщений и заканчивая удалением аккаунта.
### Авторизация посредством токеном
Пока что без подробностей. Возвращает стандартную Responce структуру с code состоянием. Всё,что не является 20X состоянием - авторизация не успешна.
```
{
	"sessia": {
		"tocken": "look example",
		"device": {
			"deviceid": "look example"
		},
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": "look example"
		}
	}
}
```
Возвращаемая структура.
```
{
    "Status": {
        "code": ,
        "message": ""
    },
    "Tocken": ""
}
```
## Данные регистрации
Для отправки регистрационной информации на сервер проходит в три этапа.
#### 1 Этап. Отправьте сведения о пользователе. /Registration
```
{
	"login": {
		"login": "look example"
	},
	"password": {
		"password": "look example"
	},
	"email": {
		"email": "look example"
	},
	"sessia": {
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": "look example"
		},
		"device": {
			"deviceid": "look example",
			"os": "look example",
			"mac": "look example",
			"hostname": "look example"
		}
	}
}
```
Обработать её в бинарный вид и отправить на сервер /.
В случае успешной регистрации, сервер вернёт структуру следующего шаблона:
```
{
    "Status": {
        "code": 0,
        "message": "look example"
    },
    "Tocken": "example"
}
```
Пример отправляемых данных:
```
{
    "login": {
        "login": "KrTv"
    },
    "password": {
        "password": "DisTet"
    },
    "email": {
        "email": "bigbro1213@inbox.ru"
    },
    "sessia": {
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82",
            "os": "Windows 10",
            "mac": "2cab5f8c00125bcc",
            "hostname": "DESKTOP-G123VB"
        }
    }
}
```
Пример ответа:
```
{
    "Status": {
        "code": 200,
        "message": "Operation done"
    },
    "Tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1"
}
```
#### 2 Этап. Установить фото профиля
Актуальная инструкция расположена в разделе [Media](http://distet.tech/VKR/server/wiki/Transferred-data#%D1%81%D0%BC%D0%B5%D0%BD%D0%B0-%D1%84%D0%BE%D1%82%D0%BE-%D0%BF%D1%80%D0%BE%D1%84%D0%B8%D0%BB%D1%8F-%D0%BD%D0%B0%D0%BF%D1%80%D0%B8%D0%BC%D0%B5%D1%80-%D0%BF%D1%80%D0%B8-%D1%80%D0%B5%D0%B3%D0%B8%D1%81%D1%82%D1%80%D0%B0%D1%86%D0%B8%D0%B8)
#### 3 Этап. Подтвердите почту пользователя. 
Актуальная инструкция расположена в разделе [Email Code](http://distet.tech/VKR/server/wiki/Transferred-data#email-code)
## Email Code 
Request - Email.
Для управления подтверждением почты, используется JSON состоящий фактически из двух изменяемых полей:
```
{
    "email":"look example",
    "code":"look example",
    "sessia": {
        "tocken": "look example",
        "device": {
            "deviceid": "look example"
        },
        "client": {
            "hashcert": "look example",
            "nameclient": "look example",
            "version": look example
        }
    }
}
```
* Изменение поля Email приведёт к попытке получения нового кода.
* Изменение поля Code приведёт к попытке активации почтового адреса.
#### Внимание, это не мультизапрос, заполенние обоих полей приведёт к выполнение только запроса на получение поля. 
Пример отправляемых данных
#### Смена кода
```
{
    "email":"example@example.ru",
    "sessia": {
        "tocken": "23db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca3?2",
        "device": {
            "deviceid": "dsa34hjrdsgth31244ghh4" example"
        },
        "client": {
            "hashcert": "23db6982caef9e9152f1a5b2589e6ca3",
            "nameclient": "RollaDieApp",
            "version": 2
        }
    }
}
```
Получаемый ответ в случае успешной высылки сообщения:
```
{
    "Status": {
        "code": 200,
        "message": "Код сменён и отправлен пользователю на почту"
    }
}
```
#### Подтверждение почты
```
{
    "code":"457213",
    "sessia": {
        "tocken": "23db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca323db6982caef9e9152f1a5b2589e6ca3?2",
        "device": {
            "deviceid": "dsa34hjrdsgth31244ghh4" example"
        },
        "client": {
            "hashcert": "23db6982caef9e9152f1a5b2589e6ca3",
            "nameclient": "RollaDieApp",
            "version": 2
        }
    }
}
```
Получаемый ответ в случае успешной высылки сообщения:
```
{
    "Status": {
        "code": 200,
        "message": "почта подтверждена"
    }
}
```
## Media
### Смена фото профиля (Например, при регистрации)
Смена фото профиля проходит в 2 действия.
1 Действие. Отправить фото на сервер. Request - Media.
```
{
    "extension":"look example",
    "bytes":"look example",
	"sessia": {
        "tocken": "look example",
        "device": {
            "deviceid": "look example"
        },
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": "look example"
		}
	}
}
```
Поле extension используется для указания расширения отправляемого файла.
Поле bytes используется для указания байтов отправляемого файла. Байты отправляются в кодировке BASE64.
Пример отправляемых данных:
```
{
    "extension":"png",
    "bytes":"c2RhYXNkYQ==",
    "sessia": {
        "tocken": "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30352d32382031323a33313a33342e3731520e986233307ad1bca887e9e707f72482967eee61c637baf6c21d1bd67021250da1da2c63041b84d8432c0ddf9b1cb9d41d8cd98f00b204e9800998ecf8427e1878e7500e60c7fd3afc92bd6d5a972ec91c03bf88be26044d8c65b618c1d35c81d83e39086f2f913fc07bfa171b5da488cd8cb0d8849e8a4fc53b17666ba912?1",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "34703b34786f3929a28fc0e0a90e9fe0",
            "nameclient": "RollAPP",
            "version": 1
        }
    }
}
```
Пример получаемых данных:
```
{
    "Status": {
        "code": 200,
        "message": "Медиа загружены"
    },
    "Key": "fa5acdf109bdee4b6f2eb886e53fc6f7"
}
```
2 Действие. Изменить фото пользователя. Request - ChangeUser.
```
{
	"imagekey": "look example",
	"sessia": {
		"tocken": "look example",
		"device": {
			"deviceid": "look example"
		},
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": look example
		}
	}
}
```
Где imageKey - ключ изображения;
В качестве ответа вы получите json следующего формата 
```
{
    "Status": {
        "code": 200,
        "message": "Запрос распознан и выполнен."
    },
    "Messages": [
        "фото профиля сменено"
    ]
}
```
### Отправка медиа
Для отправки медиа информации на сервер нужно:
* Отправить медиа на сервер. Request - Media.
```
{
    "extension":"png",
    "bytes":"MTIzMTE0MQ==",
    "sessia": {
        "tocken": "fd?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
Поле extension используется для указания расширения отправляемого файла.
Поле bytes используется для указания байтов отправляемого файла. Шифруется в BASE64.
```
{
    "Status": {
        "code": 200,
        "message": "Медиа загружены"
    },
    "Key": "fa5acdf109bdee4b6f2eb886e53fc6f7"
}
```
### Вытягивание медиа файла
Для вытягивания информации сообщений с сервера нужно:
* Отправить запрос на сервер. Request - GetMedia.
```
{
    "hash": "sa148a8cx43safsedfs",
    "sessia": {
        "tocken": "fd?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
Поле hash используется для указания к какому файлу необходимо выполнить запрос. 
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "файл сформирован."
    },
    "BytesBase64": "sa148a8cx43safsedfs",
    "Extension": "png",
    "Hash": "sa148a8cx43safsedfs"
}
```
### Вытягивание медиа контента беседы
Для вытягивания информации сообщений с сервера нужно:
* Отправить запрос на сервер. Request - GetMedia.
```
{
    "take": 48,
    "skip": 0,
    "PkDialog": "31323032332d30352d32382031343a33333a31342e393646542a29ddcd4e31a7070ac6fea1ee0bef75f4071965b1c27b381ac1db2dd0b6d2e5",
    "sessia": {
        "tocken": "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30352d32382031373a35333a32312e3237e075ce7a467d8e6a97810d08fa1dc24e63cf26d58103c0478c9be52d341d09246a2eb8020e343d17fe5021e84c5ba95ad41d8cd98f00b204e9800998ecf8427e1878e7500e60c7fd3afc92bd6d5a972ec91c03bf88be26044d8c65b618c1d35c81d83e39086f2f913fc07bfa171b5da488cd8cb0d8849e8a4fc53b17666ba912?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "34703b34786f3929a28fc0e0a90e9fe0",
            "nameclient": "RollAPP",
            "version": 1
        }
    }
}
```
Поле pkDialog используется для указания по отношению к какому диалогу необходимо выполнить запрос.
Поле Take и Skip описывают сколько за один запрос необходимо вытащить записей. Например, чтобы получить все сообщения пользователя, рекомендуется использовать цикл где skip изменяется с шагом 50, а take всегда равен 50. 
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "файл сформирован."
    },
    "BytesBase64": "sa148a8cx43safsedfs",
    "Extension": "png"
}
```
## Друзья и чёрный список
Все запросы этой подсистемы отправляются по пути /RequestFrBl. Для всех запросов используется вариантивная JSON структура, полная версия которой описана ниже:
```
{
	"id": look example,
	"description": "look example",
	"operation": "look example",
	"take": look example,
	"skip": look example,
	"sessia": {
		"tocken": "look example",
		"device": {
			"deviceid": "look example"
		},
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": look example
		}
	}
}
```
Поле ID используется для указания по отношению к кому необходимо выполнить запрос. Например, чтобы добавить пользователя в друзья, у которого ID равен 20, необходимо указать именно это число.
Поле Description используется для указание заметки при добавлении в друзья.
Поле Operation используется для указания типа запроса, перечислим допустимые запросы:
* "accept" – принять запрос в друзья
* "append" – отправить запрос в друзья 
* "remove" – удалить из друзей
* "cancel" – отозвать запрос в друзья
* "update" – обновить описание друга 
* "getFriend" – Получить список друзей
* "getIncoming" – Получить список запросов в друзья 
* "getOutcoming" – Получить список отправленных запросов в друзья
* "addBlock" – Добавить в чёрный список
* "delBlock" – Удалить из чёрного списка
* "listBlocks" – Посмотреть добавленных в чёрный список
Поле Take и Skip описывают сколько за один запрос необходимо вытащить записей. Например, чтобы получить все сообщения пользователя, рекомендуется использовать цикл где skip изменяется с шагом 100, а take всегда равен 100. 
Рассмотрим примеры отправляемых JSON.
#### ACCEPT
```
{
	"id": 1,
	"operation": "accept",
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### APPEND
```
{
	"id": 4,
	"description": "Друг со школы",
	"operation": "append",
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### REMOVE
```
{
	"id": 1,
	"operation": "remove",
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### CANCEL
```
{
	"id": 1,
	"operation": "cancel",
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### update
```
{
	"id": 4,
	"description": "Друг НЕ со школы",
	"operation": "update",
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### GETFRIEND
```
{
	"operation": "getFriend",
	"take": 50,
	"skip": 0,
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### GETINCOMING
```
{
	"operation": "getIncoming",
	"take": 50,
	"skip": 0,
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### GETOUTCOMING
```
{
	"operation": "getOutcoming",
	"take": 50,
	"skip": 0,
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### ADDBLOCK
```
{
	"id": 1,
	"operation": "addBlock",
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### DELBLOCK
```
{
	"id": 1,
	"operation": "delBlock",
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
#### LISTBLOCK
```
{
	"operation": "listBlocks",
	"take": 50,
	"skip": 0,
	"sessia": {
		"tocken": "20a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30342d32382031393a31323a35310c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0047215ee9c7d9dc229d2921a40e899ec5f20a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?1",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
		}
	}
}
```
## Вытягивание сообщений
Для вытягивания информации сообщений с сервера нужно:
* Отправить запрос на сервер. Request - GetMessages.
* Пример вытягивания всех сообщений из диалога:
```
{
    "pkDialog": "3",
    "date": "-1",
    "take": 15,
	"skip": 50,
    "sessia": {
        "tocken": "fd?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
* Пример вытягивания конкретного сообщения из диалога:
```
{
    "msgkey": "afc259dd5bc5253c16d2876f3e5532c3dfa65f7305d418e62d954e446550e784",
    "pkDialog": "31323032332d30352d32382031343a33333a31342e393646542a29ddcd4e31a7070ac6fea1ee0bef75f4071965b1c27b381ac1db2dd0b6d2e5",
    "sessia": {
        "tocken": "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30352d32382031373a35333a32312e3237e075ce7a467d8e6a97810d08fa1dc24e63cf26d58103c0478c9be52d341d09246a2eb8020e343d17fe5021e84c5ba95ad41d8cd98f00b204e9800998ecf8427e1878e7500e60c7fd3afc92bd6d5a972ec91c03bf88be26044d8c65b618c1d35c81d83e39086f2f913fc07bfa171b5da488cd8cb0d8849e8a4fc53b17666ba912?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "34703b34786f3929a28fc0e0a90e9fe0",
            "nameclient": "RollAPP",
            "version": 1
        }
    }
}
```
Поле msgkey используется для указания по отношению к какому сообщению необходимо выполнить запрос.
Поле pkDialog используется для указания по отношению к какому диалогу необходимо выполнить запрос.
Поле date используется для указания с какой даты необходимо вытягивать значения. В случае -1 вернутся 50 последних. 
Поле Take и Skip описывают сколько за один запрос необходимо вытащить записей. Например, чтобы получить все сообщения пользователя, рекомендуется использовать цикл где skip изменяется с шагом 50, а take всегда равен 50. 
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "Список сформирован."
    },
    "Msg": [
        {
            "infoMessage": {
                "msgKey": "5accfeddc9afdec25e45cca089bac5b0a4a84916a003a742e6744ae2443997a2",
                "idUser": 4,
                "date": "2023-05-14 20:30:03.08",
                "content": "Hello world",
                "updateAt": "",
                "deleteAt": "",
                "forwardedKey": "",
                "read": "",
                "important": ""
            },
            "hash": null,
            "order": null
        }
    ]
}
```
## Вытягивание удаленных сообщений
Для вытягивания информации сообщений с сервера нужно:
* Отправить запрос на сервер. Request - GetDelMessages.
* Пример вытягивания всех сообщений из диалога:
```
{
    "pkDialog": "3",
    "date": "-1",
    "take": 15,
	"skip": 50,
    "sessia": {
        "tocken": "fd?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
Поле msgkey используется для указания по отношению к какому сообщению необходимо выполнить запрос.
Поле pkDialog используется для указания по отношению к какому диалогу необходимо выполнить запрос.
Поле date используется для указания с какой даты необходимо вытягивать значения. В случае -1 вернутся 50 последних. 
Поле Take и Skip описывают сколько за один запрос необходимо вытащить записей. Например, чтобы получить все сообщения пользователя, рекомендуется использовать цикл где skip изменяется с шагом 50, а take всегда равен 50. 
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "Список сформирован."
    },
    "Msg": [
        {
            "infoMessage": {
                "msgKey": "5accfeddc9afdec25e45cca089bac5b0a4a84916a003a742e6744ae2443997a2",
                "idUser": 4,
                "date": "2023-05-14 20:30:03.08",
                "content": "Hello world",
                "updateAt": "",
                "deleteAt": "",
                "forwardedKey": "",
                "read": "",
                "important": ""
            },
            "hash": null,
            "order": null
        }
    ]
}
```
## Отправка Сообщения
Все запросы связанные с передаваемыми сообщениями клиента отправляются по пути /SendMessage. Для всех запросов используется вариативная JSON структура, пример которой описан ниже:
```
{
    "operation": "add",
    "message":{
        "infoMessage":{
            "msgkey": "de90af0d0aaceb6e4f4d77f4f69396b6e32d53d644679dbcd5aecbb0eae341fb",
            "content": "Это сообщение 3 пользователя",
            "pkDialog": "32323032332d30352d31342031373a30333a35322e3239d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
        },
          "hash": ["fa5acdf109bdee4b6f2eb886e53fc6f7"],
          "order": [1]
    },
    "sessia": {
        "tocken": "fd?3",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
Поле msgKey используется для указания по отношению к какому сообщению необходимо выполнить запрос. 
Поле content используется для указания текста отправляемого сообщения.
Поле pkDialog используется для указания по отношению к какому диалогу необходимо выполнить запрос.
Поле Operation используется для указания типа запроса, перечислим допустимые запросы:
* "add" – отправить сообщение
* "edit" – отредактировать сообщение
* "delete" – удалить сообщение
* "deleteOnlyMe" – удалить сообщение только у себя
В случае, если вы хотите добавить медиа контент к сообщению необходимо добавить следующие поля:
Поле hash используется для указания по отношению к какому медиа контенту необходимо выполнить запрос. 
Поле order используется для указания порядка отображений изображения.
Рассмотрим примеры отправляемых JSON.
#### add
```
{
    "operation": "add",
    "message":{
        "infoMessage":{
            "content": "Это сообщение 3 пользователя",
            "pkDialog": "32323032332d30352d31342031373a30333a35322e3239d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
        }
    },
    "sessia": {
        "tocken": "fd?3",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
#### Сообщение с медиа
```
{
    "operation": "add",
    "message":{
        "infoMessage":{
            "content": "Это сообщение 3 пользователя",
            "pkDialog": "32323032332d30352d31342031373a30333a35322e3239d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
        },
        "hash": ["fa5acdf109bdee4b6f2eb886e53fc6f7"],
        "order": [1]
    },
    "sessia": {
        "tocken": "fd?3",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
#### edit
```
{
    "operation": "edit",
    "message":{
        "infoMessage":{
            "msgkey": "de90af0d0aaceb6e4f4d77f4f69396b6e32d53d644679dbcd5aecbb0eae341fb",
            "content": "Это сообщение 3 пользователя",
            "pkDialog": "32323032332d30352d31342031373a30333a35322e3239d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
        }
    },
    "sessia": {
        "tocken": "fd?3",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
#### delete или deleteOnlyMe
```
{
    "operation": "delete" *Или deleteOnlyMe
    "message":{
        "infoMessage":{
            "msgkey": "de90af0d0aaceb6e4f4d77f4f69396b6e32d53d644679dbcd5aecbb0eae341fb",
            "pkDialog": "32323032332d30352d31342031373a30333a35322e3239d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
        }
    },
    "sessia": {
        "tocken": "fd?3",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
Ответ приходит по следующему образцу:
```
{
    "Status": {
        "code": 200,
        "message": "Сообщение отправлено."
    },
    "msgkey": "f262756d04f030755743b2f6dfbf983ba6fe01973fb6368be4bea6e9507b9ef9"
}
```
## Создание беседы
Для отправки запроса создания беседы на сервер необходимо отправить запрос по пути /CreateDialog. и подготовить текст в следующем виде:
```
{
    "dialog": {
        "name": "xczxczxasdasdasczxczxcxz",
        "private": "true",
        "idUsers": [3,4]
        "photo" : "99fds9fds9f9sd9f9dsf"
    },
    "sessia": {
        "tocken": "4445534b544f502d473132335642a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30352d31312031353a33373a31322e36340c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f0046a847ff66407c4ed8b96f4ad381a8aef3066636436303835393138383433643239636336643138656465303036653832a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
Поле photo используется для указания ключа фотографии диалога. 
Поле name используется для указания названия отправляемой беседы.
Поле private используется для указания приватности беседы. В случае false она не будет отображатся при поиске бесед.
Поле idUsers используется для указания по отношению к каким пользователям необходимо выполнить запрос. 
Ответ приходит по следующему образцу:
```
{
    "Status": {
        "code": 200,
        "message": "Диалог начат."
    },
    "hash": "32323032332d30352d31372032303a32353a34332e3733d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
}
```
## Изменение сведений диалога
Для того, чтобы выполнить запрос, заполните соответствующие поля по примеру и отправьте по пути **/ChangeDialog** :
```
{
    "dialog": {
        "name": "OkTested",
        "private": "0",
        "photo": "da470f8db2cace85ccb977ae8f9b84aa",
        "IdUsers": [51],
        "pkDialog": "31323032332d30352d32382031343a33333a31342e393646542a29ddcd4e31a7070ac6fea1ee0bef75f4071965b1c27b381ac1db2dd0b6d2e5"
    },
    "sessia": {
        "tocken": "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30352d32382031373a35333a32312e3237e075ce7a467d8e6a97810d08fa1dc24e63cf26d58103c0478c9be52d341d09246a2eb8020e343d17fe5021e84c5ba95ad41d8cd98f00b204e9800998ecf8427e1878e7500e60c7fd3afc92bd6d5a972ec91c03bf88be26044d8c65b618c1d35c81d83e39086f2f913fc07bfa171b5da488cd8cb0d8849e8a4fc53b17666ba912?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "34703b34786f3929a28fc0e0a90e9fe0",
            "nameclient": "RollAPP",
            "version": 1
        }
    }
}
```
Поле pkDialog используется для указания по отношению к какому диалогу необходимо выполнить запрос.
Где заполнение поля:
1) name приведёт к смене имени диалога на новое;
2) private приведёт к смене приватности диалога;
3) photo приведёт к смене фото диалога;
4) IdUsers приведёт к добавлению новых пользователей в диалог.
А заполнение всех полей приведёт к выполнению всех операций. 
В качестве ответа вы получите json следующего формата 
```
{
    "Status": {
        "code": 500,
        "message": "Данные успешно изменены."
    }
}
```
## Кик пользователя из диалога
Для изгнания пользователя из диалога нужно:
* Отправить запрос на сервер. Request - KickUser. Пример запроса:
```
{
    "PkDialog": "31323032332d30352d32382031343a33333a31342e393646542a29ddcd4e31a7070ac6fea1ee0bef75f4071965b1c27b381ac1db2dd0b6d2e5",
    "user": 51,
    "sessia": {
        "tocken": "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30352d32382031373a35333a32312e3237e075ce7a467d8e6a97810d08fa1dc24e63cf26d58103c0478c9be52d341d09246a2eb8020e343d17fe5021e84c5ba95ad41d8cd98f00b204e9800998ecf8427e1878e7500e60c7fd3afc92bd6d5a972ec91c03bf88be26044d8c65b618c1d35c81d83e39086f2f913fc07bfa171b5da488cd8cb0d8849e8a4fc53b17666ba912?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "34703b34786f3929a28fc0e0a90e9fe0",
            "nameclient": "RollAPP",
            "version": 1
        }
    }
}
```
Где
Поле pkDialog используется для указания по отношению к какому диалогу необходимо выполнить запрос.
Поле user используется для указания по отношению к какому пользователю необходимо выполнить запрос.
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "Пользователь изгнан"
    }
}
```
## Вытягивание диалогов
Для вытягивания информации диалогов с сервера нужно:
* Отправить запрос на сервер. Request - GetDialogs.
```
{
    "date": "-1",
    "take": 15,
	"skip": 50,
    "sessia": {
        "tocken": "fd?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "08520331451320871a4e77a60d70dbe2",
            "nameclient": "RollaDieApp",
            "version": 1
        }
    }
}
```
Поле date используется для указания с какой даты необходимо вытягивать значения. В случае -1 вернутся 50 последних. 
Поле Take и Skip описывают сколько за один запрос необходимо вытащить записей. Например, чтобы получить все сообщения пользователя, рекомендуется использовать цикл где skip изменяется с шагом 50, а take всегда равен 50. 
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "Список сформирован."
    },
    "Dialogs": [
        {
            "hash": "32323032332d30352d31372032303a32353a34332e3733d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
            "name": "Это диалог",
            "photo": ""
        }
    ]
}
```
## Вытягивание данных пользователя
Для вытягивания информации данных пользователей с сервера нужно:
* Отправить запрос на сервер. Request - GetUser. Пример запроса:
```
{
    "id": 1,
    "sessia": {
        "tocken": "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a323032332d30352d32382031373a35333a32312e3237e075ce7a467d8e6a97810d08fa1dc24e63cf26d58103c0478c9be52d341d09246a2eb8020e343d17fe5021e84c5ba95ad41d8cd98f00b204e9800998ecf8427e1878e7500e60c7fd3afc92bd6d5a972ec91c03bf88be26044d8c65b618c1d35c81d83e39086f2f913fc07bfa171b5da488cd8cb0d8849e8a4fc53b17666ba912?2",
        "device": {
            "deviceid": "0fcd6085918843d29cc6d18ede006e82"
        },
        "client": {
            "hashcert": "34703b34786f3929a28fc0e0a90e9fe0",
            "nameclient": "RollAPP",
            "version": 1
        }
    }
}
```
Поле id используется для указания пользователя, данные которого необходимо получить. 
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "Данные пользователя сформированы."
    },
    "Photo": "",
    "Nickname": "Login"
    "Id": 1
}
```
## Вытягивание участников диалога
Для вытягивания информации участников диалога с сервера нужно:
* Отправить запрос на сервер. Request - GetUsersDialog. Пример запроса:
```
{
    "pkDialog": "2",
    "sessia": {
        "client": {
            "hashcert": "990609802ba87c27f0c09195e57385d9",
            "nameclient": "Rolladie",
            "version": 1
        },
        "device": {
            "deviceid": "362c3e8a59124fdf8f7481839e6ba4aa"
        },
        "tocken": "4b7254764d6437c1e17dcf310a72a98a5903285ee97f77a2dd505496c9ecd5f093530c91323032332d30362d30342031363a31343a32342e36309353f41348ec8e3bba297f55ee083c99a9dfb910c0222076e8d0e49c032643b91344b5e9b06d207d0b111f72aae244d34498efa6a6892c8d38503f9012d6bac94503f56278516b9d602bb88503f79849b1e49726561f452e52c9b18f48282c90a47b5470c6cbcbbc8d6a83188e94bbdc6467252dd8ca9d0fe3510d83d75ac0c3?1"
    }

}
```
Поле pkDialog используется для указания диалога, участников которого необходимо получить. 
Пример возвращаемого значения:
```
{
    "Status": {
        "code": 200,
        "message": "Данные пользователей переписки сформированы."
    },
    "UsersId": [
        2,
        1
    ]
}
```
## Поиск
Поиск разделяется на две системы. Поиск внутри диалога и поиск во всех диалогах. Рассмотрим то, как работает последнее. Для поиска необходимо использовать вариативную JSON структуру:
```
{
	"dialogs": "look example",
	"search": "look example",
	"take": look example,
	"takedialogs": look example,
	"skip": look example,
	"skipdialogs": look example,
	"operation": "look example",
	"sessia": {
		"tocken": look example,
		"device": {
			"deviceid": "look example"
		},
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": look example
		}
	}
}
```
Подготовим запрос для поиска во всех диалогах JSON:
```
{
	"search": "тестовые данные",
	"takedialogs": 50,
	"skipdialogs": 0,
	"operation": "everyWhere",
	"sessia": {
		"tocken": "fd?2",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
			"hashcert": "08520331451320871a4e77a60d70dbe2",
			"nameclient": "RollaDieApp",
			"version": 1
		}
	}
}
```
Возвращаемое значение в таком случае будет примерно таким:
```
{
    "Status": {
        "code": 500,
        "message": "Внутренняя ошибка при поиске"
    },
    "Dialogs": [
        "32323032332d30352d32312031383a33383a30322e38334654a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
        "3",
        "32323032332d30352d31372032303a32353a34332e3733d0add182d0be20d0b4d0b8d0b0d0bbd0bed0b3a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"
    ]
}
```
После чего, необходимо взять все ключи диалогов и отправить их следующим запросом:
```
{
	"dialogs": "32323032332d30352d32312031383a33383a30322e38334654a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
	"search": "тестовые данные",
	"take": 50,
	"skip": 0,
	"operation": "inDialogs",
	"sessia": {
		"tocken": "fd?2",
		"device": {
			"deviceid": "0fcd6085918843d29cc6d18ede006e82"
		},
		"client": {
			"hashcert": "08520331451320871a4e77a60d70dbe2",
			"nameclient": "RollaDieApp",
			"version": 1
		}
	}
}
```
После отправки запроса, вернётся JSON с массивом ключей диалогов, где содержатся ключи сообщений:
```
{
    "Status": {
        "code": 200,
        "message": "Список сообщений сформирован."
    },
    "Messages": [
        "de90af0d0aaceb6e4f4d77f4f69396b6e32d53d644679dbcd5aecbb0eae341fb",
        "2c48cad9da9c226c96b625e41855995f65f732f28f240bf7c5158353623bc1be",
        "5a8b49d2705ac7442dec1c2d5f5b4dbc4b2b4f679f51ac8e83c0b130e4c56105",
        "32cb14cb0fbe6d27ab06c88819c6d7d627f394453dd7e19fac62daaaf85b4d97",
        "693566719b98bf14cd2bcd741913f6d353d5a493d61ec3d6e465cf4bc98d5afc",
        "9b03a53f0163045bff4731e9bf6c1e3c6880b23ddf7ec13af5dd23a97b69d793"
    ]
}
```
Чтобы получить конкретную информацию о сообщениях, обратиесь в подраздел "загрузка сообщения по ключу"
## Изменение пользовательских сведений
Для того, чтобы выполнить запрос, заполните соответствующие поля в шаблоне и отправьте по пути **/ChangeUser** :
```
{
	"nickname": "look example",
	"oldPass": "look example",
	"newPass": "look example",
	"imagekey": "look example",
	"email": "look example",
	"sessia": {
		"tocken": "look example",
		"device": {
			"deviceid": "look example"
		},
		"client": {
			"hashcert": "look example",
			"nameclient": "look example",
			"version": look example
		}
	}
}
```
Где заполнение поля:
1) nickname приведёт к смене имени пользователя на новое;
2) oldPass и newPass приведёт к смене пароля;
3) imageKey приведёт к смене фото профиля;
4) email приведёт ни к чему.
А заполнение всех полей приведёт к выполнение всех операций. 
В качестве ответа вы получите json следующего формата 
```
{
    "Status": {
        "code": 200,
        "message": "Список сообщений сформирован."
    },
    "Messages": [
        "никнейм сменён",
        "пароль сменён",
        "фото профиля сменено"
    ]
}
```
## Выход из системы
Отправьте любую JSON из вышеприведённых, содержащее поле "sessia" по пути **/disable**