def is_safe(board, row, col, n):
    for i in range(row):
        if board[i] == col:
            return False
    
    for i in range(row):
        if abs(board[i] - col) == abs(i - row):
            return False
    
    return True

def solve_n_queens(board, row, n):
    if row == n:
        print_board(board, n)
        return True
    
    for col in range(n):
        if is_safe(board, row, col, n):
            board[row] = col
            if solve_n_queens(board, row + 1, n):
                return True
            board[row] = -1 
    
    return False

def print_board(board, n):
    for row in range(n):
        line = ['.'] * n
        line[board[row]] = 'Q'
        print(' '.join(line))
    print()

n = 8
board = [-1] * n

solve_n_queens(board, 0, n)
