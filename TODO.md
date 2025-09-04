# api

1. Запилить максимально просто быстро api запросы, данные хранить в map
Использовать fiber
    1. `POST` `/user` req: {name, age, ...}; resp id(uuid) 201, 400, 500
    2. `GET` `/user/:id`; resp {id, name, age, ...} 200, 400, 500
    3. `DELETE` `/user/:id`; resp 200, 400, 500
    4. `PUT` `/user` req: {name, age, ...}; resp id 200, 400, 500

1. слайс vs массив, что под капотом, что происходит при append (cap)
2. map, бакеты, коллизии, эвакуации, swiss table в новой версии
3. interface, что подкапотом, solid (чек примеры в go), опп в го примеры чек

https://www.youtube.com/@Skills_mentor/videos
