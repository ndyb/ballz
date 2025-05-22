<script lang="ts">
	import { onMount } from 'svelte';

	// Game configuration
	const shapes = ['circle', 'square', 'triangle', 'hexagon'];
	const colors = ['#FF5252', '#4CAF50', '#2196F3', '#FFC107'];

	let gameStarted = false;
	let score = 0;
	let timeLeft = 30;
	let timer: any;
	let gameBoard: any = [];
	let targetShape: any = null;
	let targetColor: any = null;

	// Create a new game board with random shapes and colors
	function createBoard() {
		gameBoard = [];
		for (let i = 0; i < 16; i++) {
			gameBoard.push({
				id: i,
				shape: shapes[Math.floor(Math.random() * shapes.length)],
				color: colors[Math.floor(Math.random() * colors.length)]
			});
		}
		setNewTarget();
	}

	// Set a new target shape and color
	function setNewTarget() {
		targetShape = shapes[Math.floor(Math.random() * shapes.length)];
		targetColor = colors[Math.floor(Math.random() * colors.length)];
	}

	// Handle click on a shape
	function handleClick(item: any) {
		if (!gameStarted) return;

		if (item.shape === targetShape && item.color === targetColor) {
			score += 10;
			// Replace the clicked shape
			const index = gameBoard.findIndex((shape: any) => shape.id === item.id);

			// Generate a new shape and color that are different from the target
			let newShape, newColor;
			do {
				newShape = shapes[Math.floor(Math.random() * shapes.length)];
				newColor = colors[Math.floor(Math.random() * colors.length)];
			} while (newShape === targetShape && newColor === targetColor);

			gameBoard[index] = {
				id: item.id,
				shape: newShape,
				color: newColor
			};

			setNewTarget();
		} else {
			score = Math.max(0, score - 5);
		}

		gameBoard = [...gameBoard]; // Trigger reactivity
	}

	// Start the game
	function startGame() {
		gameStarted = true;
		score = 0;
		timeLeft = 30;
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
