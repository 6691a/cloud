package dns

type HandlerFunc func(task Task, dns DNS) Response

var handlers = map[Method]HandlerFunc{
	Method(Get):    handlerGet,
	Method(List):   handlerList,
	Method(Create): handlerCreate,
	Method(Patch):  handlerPatch,
	Method(Delete): handlerDelete,
}

func handlerGet(task Task, dns DNS) Response {
	record, err := dns.Get(task.Request.Record.SubDomain, task.Request.Record.Type)
	return Response{
		Error:   err,
		Records: []Record{record},
	}
}

func handlerList(task Task, dns DNS) Response {
	records, err := dns.List()
	return Response{
		Error:   err,
		Records: records,
	}
}

func handlerCreate(task Task, dns DNS) Response {
	record, err := dns.Create(task.Request.Record)
	return Response{
		Error:   err,
		Records: []Record{record},
	}
}

func handlerPatch(task Task, dns DNS) Response {
	record, err := dns.Patch(task.Request.Record.SubDomain, task.Request.Record.Type, task.Request.Record)
	return Response{
		Error:   err,
		Records: []Record{record},
	}
}

func handlerDelete(task Task, dns DNS) Response {
	err := dns.Delete(task.Request.Record.SubDomain, task.Request.Record.Type)
	return Response{
		Error:   err,
		Records: []Record{},
	}
}
