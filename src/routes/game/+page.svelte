<script lang="ts">
	import { onMount } from 'svelte';

	// Game configuration
	const shapes = ['circle', 'square', 'triangle', 'hexagon'];
	const colors = ['#FF5252', '#4CAF50', '#2196F3', '#FFC107'];

	let gameStarted = false;
	let score = 0;
	let timeLeft = 60;
	let timer: any;
	let gameBoard: any = [];
	let targetShape: string = '';
	let targetColor: string = '';

	// Create a new game board with random shapes and colors
	function createBoard() {
		gameBoard = [];
		for (let i = 0; i < 16; i++) {
			gameBoard.push({
				id: i,
				shape: shapes[Math.floor(Math.random() * shapes.length)],
				color: colors[Math.floor(Math.random() * colors.length)],
				row: Math.floor(i / 4),
				col: i % 4
			});
		}
	}

	// Handle click on a shape
	function handleClick(item: any) {
		if (!gameStarted) return;

		// Find adjacent matching shapes
		const matches = findMatches(item);

		if (matches.length > 0) {
			// Add the clicked item to matches
			matches.push(item);

			// Award points
			score += matches.length * 10;

			// Remove matches from board
			removeMatches(matches);

			// Make pieces fall down
			applyGravity();

			// Check if game is won
			if (gameBoard.length === 0) {
				clearInterval(timer);
				gameStarted = false;
				// Additional win logic could go here
			}
		}
	}

	// Find adjacent shapes that match the clicked shape
	function findMatches(item: any) {
		const matches: any[] = [];
		const clickedShape = item.shape;
		const clickedColor = item.color;

		// Check adjacent positions (left, right, above, below)
		const directions = [
			{ row: 0, col: -1 }, // left
			{ row: 0, col: 1 }, // right
			{ row: -1, col: 0 }, // above
			{ row: 1, col: 0 } // below
		];

		directions.forEach((dir) => {
			const adjacentRow = item.row + dir.row;
			const adjacentCol = item.col + dir.col;

			// Find the adjacent piece if it exists
			const adjacentPiece = gameBoard.find(
				(p: any) => p.row === adjacentRow && p.col === adjacentCol
			);

			// If there's an adjacent piece and it matches the shape and color, add to matches
			if (
				adjacentPiece &&
				adjacentPiece.shape === clickedShape &&
				adjacentPiece.color === clickedColor
			) {
				matches.push(adjacentPiece);
			}
		});

		return matches;
	}

	// Remove matched shapes from the board
	function removeMatches(matches: any[]) {
		const matchIds = matches.map((m) => m.id);
		gameBoard = gameBoard.filter((item: any) => !matchIds.includes(item.id));
	}

	// Make pieces fall down to fill empty spaces
	function applyGravity() {
		// Create a 4x4 grid representation
		const grid: any[][] = Array(4)
			.fill(null)
			.map(() => Array(4).fill(null));

		// Place pieces in the grid
		gameBoard.forEach((piece: any) => {
			grid[piece.row][piece.col] = piece;
		});

		// Make pieces fall down column by column
		for (let col = 0; col < 4; col++) {
			let emptyRow = 3;

			// Start from bottom and move up
			for (let row = 3; row >= 0; row--) {
				if (grid[row][col] !== null) {
					// If the piece isn't already at the bottom, move it down
					if (row !== emptyRow) {
						const piece = grid[row][col];
						piece.row = emptyRow;
						grid[emptyRow][col] = piece;
						grid[row][col] = null;
					}
					emptyRow--;
				}
			}
		}

		// Update gameBoard with new positions
		gameBoard = [];
		for (let row = 0; row < 4; row++) {
			for (let col = 0; col < 4; col++) {
				if (grid[row][col] !== null) {
					gameBoard.push(grid[row][col]);
				}
			}
		}
	}

	// Start the game
	function startGame() {
		gameStarted = true;
		score = 0;
		timeLeft = 60;
		createBoard();

		timer = setInterval(() => {
			timeLeft--;
			if (timeLeft <= 0) {
				clearInterval(timer);
				gameStarted = false;
			}
		}, 1000);
	}

	// Clean up on unmount
	onMount(() => {
		return () => {
			if (timer) clearInterval(timer);
		};
	});
</script>

<div class="game-container">
	<h1>Shape Matcher</h1>

	{#if !gameStarted}
		<div class="start-screen">
			<button on:click={startGame}>Start Game</button>
			{#if score > 0}
				<p>Game Over! Your score: {score}</p>
			{/if}
		</div>
	{:else}
		<div class="game-info">
			<div class="target">
				<h3>Find this shape:</h3>
				<div class="shape {targetShape}" style="background-color: {targetColor}"></div>
			</div>
			<div class="stats">
				<p>Score: {score}</p>
				<p>Time: {timeLeft}s</p>
			</div>
		</div>

		<div class="game-board">
			{#each gameBoard as item (item.id)}
				<div class="shape-container" on:click={() => handleClick(item)}>
					<div class="shape {item.shape}" style="background-color: {item.color}"></div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.game-container {
		max-width: 600px;
		margin: 0 auto;
		padding: 20px;
		text-align: center;
	}

	.start-screen {
		display: flex;
		flex-direction: column;
		align-items: center;
		margin-top: 50px;
	}

	button {
		padding: 10px 20px;
		font-size: 18px;
		background-color: #4caf50;
		color: white;
		border: none;
		border-radius: 5px;
		cursor: pointer;
	}

	.game-board {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		grid-gap: 10px;
		margin-top: 20px;
	}

	.shape-container {
		aspect-ratio: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: #f0f0f0;
		border-radius: 5px;
		cursor: pointer;
		transition: transform 0.2s;
	}

	.shape-container:hover {
		transform: scale(1.05);
	}

	.shape {
		width: 70%;
		height: 70%;
	}

	.circle {
		border-radius: 50%;
	}

	.square {
		border-radius: 0;
	}

	.triangle {
		clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
	}

	.hexagon {
		clip-path: polygon(25% 0%, 75% 0%, 100% 50%, 75% 100%, 25% 100%, 0% 50%);
	}

	.game-info {
		display: flex;
		justify-content: space-between;
		margin-bottom: 20px;
	}

	.target {
		text-align: left;
	}

	.target .shape {
		width: 60px;
		height: 60px;
	}

	.stats {
		text-align: right;
	}
</style>
