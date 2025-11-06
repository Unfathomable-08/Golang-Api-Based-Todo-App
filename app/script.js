const todoInput = document.getElementById('todoInput');
const addBtn = document.getElementById('addBtn');
const todoList = document.getElementById('todoList');
const stats = document.getElementById('stats');

let todos = [];
let editingId = null;

// Load all todos
async function loadTodos() {
    try {
        todoList.innerHTML = '<div class="empty-state">Loading...</div>';
        const res = await fetch("http://localhost:8080/todo");
        if (!res.ok) throw new Error('Failed to fetch');
        todos = await res.json();
        renderTodos();
        updateStats();
    } catch (err) {
        todoList.innerHTML = `<div class="empty-state">⚠️ Can't connect to server<br><small>Make sure backend is running on port 8080</small></div>`;
        stats.textContent = 'Offline';
    }
}

// Render todos
function renderTodos() {
    if (todos.length === 0) {
        todoList.innerHTML = `
            <div class="empty-state">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
                            d="M9 12h6m-6-4h6m-6 8h6m-2 4H7a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v12a2 2 0 01-2 2h-4z"></path>
                </svg>
                <p>No tasks yet. Add one above!</p>
            </div>`;
        return;
    }

    todoList.innerHTML = todos.map(todo => `
        <div class="todo-item" data-id="${todo.id}">
            <div class="checkbox ${todo.completed ? 'checked' : ''}" 
                    onclick="toggleComplete(${todo.id})"></div>
            <div class="todo-text ${todo.completed ? 'completed' : ''}"
                    ondblclick="startEdit(${todo.id}, this)">
                ${escapeHtml(todo.title)}
            </div>
            <div class="actions">
                <button class="action-btn edit-btn" onclick="startEdit(${todo.id})">
                    <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                    </svg>
                </button>
                <button class="action-btn delete-btn" onclick="deleteTodo(${todo.id})">
                    <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                </button>
            </div>
        </div>
    `).join('');
}

// Add new todo
async function addTodo() {
    const title = todoInput.value.trim();
    if (!title) return;

    addBtn.disabled = true;
    addBtn.textContent = 'Adding...';

    try {
        const res = await fetch("http://localhost:8080/todo", {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title, completed: false })
        });

        if (!res.ok) throw new Error();
        const newTodo = await res.json();
        todos.push(newTodo);
        todoInput.value = '';
        renderTodos();
        updateStats();
    } catch (err) {
        alert('Failed to add task. Is backend running?');
    } finally {
        addBtn.disabled = false;
        addBtn.innerHTML = `
            <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
            Add
        `;
    }
}

// Toggle complete
async function toggleComplete(id) {
    const todo = todos.find(t => t.id === id);
    if (!todo) return;

    try {
        const res = await fetch(`http://localhost:8080/todo/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ...todo, completed: !todo.completed })
        });

        if (!res.ok) throw new Error();
        todo.completed = !todo.completed;
        renderTodos();
        updateStats();
    } catch (err) {
        alert('Failed to update');
    }
}

// Delete todo
async function deleteTodo(id) {
    if (!confirm('Delete this task?')) return;

    try {
        const res = await fetch(`http://localhost:8080/todo/${id}`, { method: 'DELETE' });
        if (!res.ok) throw new Error();

        todos = todos.filter(t => t.id !== id);
        renderTodos();
        updateStats();
    } catch (err) {
        alert('Failed to delete');
    }
}

// Edit todo
function startEdit(id, element) {
    if (editingId === id) return;
    editingId = id;

    const todo = todos.find(t => t.id === id);
    const input = document.createElement('input');
    input.type = 'text';
    input.value = todo.title;
    input.className = 'todo-edit';
    input.style.cssText = `
        flex: 1; padding: 0.5rem; border: 2px solid var(--primary); 
        border-radius: 8px; font-size: 1rem; outline: none;
    `;

    const saveEdit = async () => {
        const newTitle = input.value.trim();
        if (!newTitle || newTitle === todo.title) {
            cancelEdit();
            return;
        }

        try {
            const res = await fetch(`"http://localhost:8080/todo"/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ ...todo, title: newTitle })
            });

            if (!res.ok) throw new Error();
            todo.title = newTitle;
            renderTodos();
            updateStats();
        } catch (err) {
            alert('Failed to update');
            cancelEdit();
        }
    };

    const cancelEdit = () => {
        element.innerHTML = escapeHtml(todo.title);
        editingId = null;
    };

    input.addEventListener('blur', saveEdit);
    input.addEventListener('keydown', e => {
        if (e.key === 'Enter') saveEdit();
        if (e.key === 'Escape') cancelEdit();
    });

    element.innerHTML = '';
    element.appendChild(input);
    input.focus();
    input.select();
}

// Update stats
function updateStats() {
    const total = todos.length;
    const completed = todos.filter(t => t.completed).length;
    const pending = total - completed;

    stats.textContent = 
        total === 0 
            ? 'No tasks' 
            : `${pending} pending • ${completed} done • ${total} total`;
}

// Escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Event listeners
addBtn.addEventListener('click', addTodo);
todoInput.addEventListener('keydown', e => {
    if (e.key === 'Enter') addTodo();
});

// Initialize
loadTodos();