package organism

func NewOrganism() *Organism {
	o := &Organism{}
	o.core = o.GrowLimb("core")
	return o
}

// Organism отвечает за жизнь вашего микросервиса. Любая часть вашего микросервиса, от которой зависит его "жизнь"
// считается конечностью, например клиентский http-сервер. Организм не может считаться готовым к работе, если хотя бы
// одна из конечностей не готова. Так же организм считается мертвым, если хотя бы одна из его конечностей отмерла.
//
// Например, у вас есть http-сервер, grpc-сервер, оба эти сервера считаются конечностями и пока они оба не будут
// готовы, организм не вернет ОК на readiness пробу. Если хотя бы одна из конечностей умрет (завершится один из
// серверов), то весь организм будет считать мертвым (liveness проба не вернет ОК), и k8s увидит это в течении N секунд,
// прибьет мертвеца и развернет новый микросервис
type Organism struct {
	limbs []*Limb
	core  *Limb
}

func (o *Organism) GrowLimb(name string) *Limb {
	limb := newLimb(name)
	o.limbs = append(o.limbs, limb)

	return limb
}

func (o *Organism) IsReady() bool {
	for _, l := range o.limbs {
		if !l.IsReady() {
			return false
		}
	}

	return true
}

func (o *Organism) IsAlive() bool {
	for _, l := range o.limbs {
		if !l.IsAlive() {
			return false
		}
	}

	return true
}

func (o *Organism) DeadLimbs() []*Limb {
	var deadLimbs []*Limb
	for _, limb := range o.limbs {
		if !limb.IsAlive() {
			deadLimbs = append(deadLimbs, limb)
		}
	}

	return deadLimbs
}

func (o *Organism) NotReadyLimbs() []*Limb {
	var notReadyLimbs []*Limb
	for _, limb := range o.limbs {
		if !limb.IsReady() {
			notReadyLimbs = append(notReadyLimbs, limb)
		}
	}

	return notReadyLimbs
}

func (o *Organism) Ready() {
	o.core.Ready()
}

func (o *Organism) Die() {
	o.core.Die()
}

func newLimb(name string) *Limb {
	return &Limb{name: name, isAlive: true}
}

type Limb struct {
	name    string
	isReady bool
	isAlive bool
}

func (l *Limb) Ready() {
	l.isReady = true
}

func (l *Limb) Die() {
	l.isAlive = false
}

func (l *Limb) IsReady() bool {
	return l.isReady
}

func (l *Limb) IsAlive() bool {
	return l.isAlive
}

func (l *Limb) Name() string {
	return l.name
}
