/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
package addTwo

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	resHead := *ListNode
	resNow := *ListNode
	resHead = New(ListNode)
	resNow = resHead
	resNow.Val = l1.Val + l2.Val

	for {
		if l1.Next == nil && l2.Next == nil {
			return resHead
		} else if l1 == nil {
			for l2 != nil {
				resNow.Next = New(ListNode)
				resNow = resNow.Next
				resNow.Val = l2.Val
				l2 = l2.Next
			}
			return resHead
		} else if l2 == nil {
			for l1 != nil {
				resNow.Next = New(ListNode)
				resNow = resNow.Next
				resNow.Val = l1.Val
				l1 = l1.Next
			}
			return resHead
		} else {
			resNow.Next = New(ListNode)
			resNow = resNow.Next
			resNow.Val = l1.Val + l2.Val
			l1 = l1.Next
			l2 = l2.Next
		}
	}
	return nil

}
