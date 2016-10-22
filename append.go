package prolog

/*
%% append(list,list,list)
%% append(+,+,+) is semidet.
%% append(+,+,-) is det.
%% append(+,-,+) is semidet.
%% append(+,-,-) is det.
%% append(-,+,+) is multi.
%% append(-,+,-) is multi.
%% append(-,-,+) is multi.
%% append(-,-,-) is multi.
append(Front,Back,Whole) :- % simplified Prolog only allows variables in clause heads
	unify(Front,[]),  % simplified Prolog does unification with unify/2
	unify(Back,Whole).
append(Front,Back,Whole) :-
	cons(X,RestFront,A),  % simplified Prolog does lists with cons/3
	unify(Front,A),
	cons(X,RestWhole,B),
	unify(Whole,B),
	append(RestFront,Back,RestWhole).
*/

var empty = NewAtom("[]")

type PredicateId int

// a constant for each known predicate
const (
	predAppend3 PredicateId = iota
)

func append3_a(Front, Back, Whole Term) Goal {
	return conjunction(
		unify(Front, empty),
		unify(Back, Whole),
	)
}

func append3_b(Front, Back, Whole Term) Goal {
	var A, B, RestFront, RestWhole, X Variable
	return conjunction(
		cons3(X, RestFront, A),
		unify(Front, A),
		cons3(X, RestWhole, B),
		unify(Whole, B),
		append3(RestFront, Back, RestWhole),
	)
}

func append3(Front, Back, Whole Term) Goal {
	return &PredicateGoal{
		Predicate: predAppend3,
		Args:      []Term{Front, Back, Whole},
	}
}
