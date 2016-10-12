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

var empty = term.NewAtom("[]")

func append3_iii_a_semidet(Front,Back,Whole term.Term) (success bool) {
	for {
		if !Unify(Front,empty) {
			break
		}
		if !Unify(Back,Whole) {
			break
		}
		success = true
	}
	if !success {
		Unwind(Front)
		Unwind(Back)
		Unwind(Whole)
	}
}

func append3_iii_b_semidet(Front,Back,Whole term.Term) (success bool) {
	var A, B, RestFront, RestWhole, X term.Variable
	for {
		cons3_ooo_a_det(X,RestFront,A)  // det func has no return value
		if !Unify(Front,A) {
			break
		}
		cons3_ooo_a_det(X,RestWhole,B)  // det func has no return value
		if !Unify(Whole,B) {
			break
		}
		
		// append(RestFront,Back,RestWhole)
		if !append3_iii_semidet(RestFront,Back,RestWhole) {  // call the predicate, not the clause
			break
		}
		success = true
	}
	if !success {
		Unwind(Front)
		Unwind(Back)
		Unwind(Whole)
	}
}