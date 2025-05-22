<script lang="ts">
	import { onMount } from 'svelte';

	// Define interfaces for our score data
	interface Score {
		player: string;
		score: number;
		timestamp: string;
	}

	// Data states
	let allTimeHigh: Score | null = null;
	let dailyHigh: Score | null = null;
	let recentScores: Score[] = [];
	let loading = true;
	let error: string | null = null;

	// Format date helper function
	function formatDate(timestamp: string): string {
		const date = new Date(timestamp);
		const today = new Date();
		const yesterday = new Date(today);
		yesterday.setDate(yesterday.getDate() - 1);

		if (date.toDateString() === today.toDateString()) {
			return `Today at ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
		} else if (date.toDateString() === yesterday.toDateString()) {
			return `Yesterday at ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
		} else {
			return date.toLocaleDateString('en-US', {
				month: 'short',
				day: 'numeric',
				year: 'numeric'
			});
		}
	}

	// Fetch data from API
	async function fetchScores() {
		try {
			loading = true;
			const response = await fetch(
				'https://nimble-gaufre-8a88fc.netlify.app/.netlify/functions/api/scores'
			);

			if (!response.ok) {
				throw new Error(`API error: ${response.status}`);
			}

			const data = await response.json();
			// Extract scores from the response and map to our Score interface
			const scores: Score[] = data.scores.map((item: any) => ({
				player: 'Anonymous', // Adding default player name since it's not in the API response
				score: item.score,
				timestamp: item.created
			}));

			console.log('Fetched scores:', scores);

			// Find all-time high score
			allTimeHigh = [...scores].sort((a, b) => b.score - a.score)[0] || null;

			// Find today's high score
			const today = new Date().toISOString().split('T')[0];
			const todayScores = scores.filter((score) =>
				new Date(score.timestamp).toISOString().startsWith(today)
			);
			dailyHigh =
				todayScores.length > 0 ? [...todayScores].sort((a, b) => b.score - a.score)[0] : null;

			// Get recent scores (sorted by timestamp, newest first)
			recentScores = [...scores]
				.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
				.slice(0, 10);
		} catch (err) {
			console.error('Error fetching scores:', err);
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			loading = false;
		}
	}

	// Fetch scores when component mounts
	onMount(fetchScores);
</script>

<div class="container">
	<!-- Page title -->
	<header>
		<h1>🏆 Ballz Leaderboard</h1>
		<div class="header-underline"></div>
	</header>

	{#if error}
		<div class="error-message">
			<p>Error loading scores: {error}</p>
		</div>
	{/if}

	{#if loading}
		<div class="loader-container">
			<div class="loader"></div>
		</div>
	{:else}
		<!-- Highest scores section -->
		<div class="top-scores">
			<!-- All-time high score -->
			<div class="score-card all-time">
				<h2>All-Time High</h2>
				{#if allTimeHigh}
					<p class="score-value">{allTimeHigh.score.toLocaleString()}</p>
					<p class="player-name">{allTimeHigh.player}</p>
					<p class="timestamp">{formatDate(allTimeHigh.timestamp)}</p>
				{:else}
					<p class="no-data">No data available</p>
				{/if}
			</div>

			<!-- Daily high score -->
			<div class="score-card daily">
				<h2>Today's High</h2>
				{#if dailyHigh}
					<p class="score-value">{dailyHigh.score.toLocaleString()}</p>
					<p class="player-name">{dailyHigh.player}</p>
					<p class="timestamp">{formatDate(dailyHigh.timestamp)}</p>
				{:else}
					<p class="no-data">No scores today</p>
				{/if}
			</div>
		</div>

		<!-- Recent scores table -->
		<div class="scores-table-container">
			<h2>Recent Scores</h2>
			<table>
				<thead>
					<tr>
						<th>Rank</th>
						<th>Player</th>
						<th>Score</th>
						<th>Date</th>
					</tr>
				</thead>
				<tbody>
					{#each recentScores as score, i}
						<tr>
							<td>{i + 1}</td>
							<td>{score.player}</td>
							<td class="score-cell">{score.score.toLocaleString()}</td>
							<td>{formatDate(score.timestamp)}</td>
						</tr>
					{:else}
						<tr class="empty-row">
							<td colspan="4">No recent scores available</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<style>
	.container {
		max-width: 1000px;
		margin: 0 auto;
		padding: 2rem 1rem;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}

	header {
		margin-bottom: 3rem;
	}

	header h1 {
		font-size: 3rem;
		text-align: center;
		background: linear-gradient(45deg, #4a6cf7, #2dd4bf);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		margin-bottom: 0.5rem;
		font-weight: 800;
		letter-spacing: -1px;
	}

	.header-underline {
		height: 4px;
		width: 120px;
		margin: 0 auto;
		background: linear-gradient(90deg, #4a6cf7, #2dd4bf);
		border-radius: 2px;
	}

	.error-message {
		background: linear-gradient(to right, #ffccd5, #ffd6e0);
		border-left: 4px solid #ff4d6d;
		color: #c9184a;
		padding: 1rem;
		margin-bottom: 2rem;
		border-radius: 0 4px 4px 0;
	}

	.loader-container {
		display: flex;
		justify-content: center;
		margin: 3rem 0;
	}

	.loader {
		width: 50px;
		height: 50px;
		border: 3px solid rgba(74, 108, 247, 0.3);
		border-top-color: #4a6cf7;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	.top-scores {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 2rem;
		margin-bottom: 3rem;
	}

	.score-card {
		padding: 2rem;
		border-radius: 12px;
		text-align: center;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
		transition:
			transform 0.2s,
			box-shadow 0.2s;
	}

	.score-card:hover {
		transform: translateY(-5px);
		box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
	}

	.all-time {
		background: linear-gradient(135deg, #e2f0ff, #f0f8ff);
		border-left: 4px solid #4a6cf7;
	}

	.daily {
		background: linear-gradient(135deg, #e2f8f0, #f0fffa);
		border-left: 4px solid #2dd4bf;
	}

	.score-card h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin-bottom: 1rem;
		color: #333;
	}

	.score-value {
		font-size: 3.5rem;
		font-weight: 800;
		line-height: 1;
		margin-bottom: 0.75rem;
		background: linear-gradient(45deg, #4a6cf7, #2dd4bf);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.all-time .score-value {
		font-size: 4rem;
	}

	.player-name {
		font-size: 1.1rem;
		font-weight: 500;
		color: #444;
		margin-bottom: 0.25rem;
	}

	.timestamp {
		font-size: 0.875rem;
		color: #777;
	}

	.no-data {
		font-size: 1.25rem;
		color: #999;
		margin-top: 1rem;
	}

	.scores-table-container {
		background: white;
		border-radius: 12px;
		box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
		overflow: hidden;
	}

	.scores-table-container h2 {
		text-align: center;
		padding: 1.5rem;
		margin: 0;
		font-size: 1.5rem;
		color: #333;
		border-bottom: 1px solid #f1f1f1;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	thead {
		background: linear-gradient(45deg, #f8f9fe, #f1f4fd);
	}

	th {
		padding: 1rem;
		text-align: left;
		font-weight: 600;
		color: #555;
		border-bottom: 2px solid #edf2f7;
	}

	td {
		padding: 1rem;
		border-bottom: 1px solid #edf2f7;
		color: #444;
	}

	tbody tr:hover {
		background-color: #f9fafb;
	}

	.score-cell {
		font-weight: 600;
		color: #4a6cf7;
	}

	.empty-row td {
		text-align: center;
		color: #999;
		padding: 2rem;
	}

	@media (max-width: 768px) {
		.top-scores {
			grid-template-columns: 1fr;
			gap: 1.5rem;
		}

		.score-value,
		.all-time .score-value {
			font-size: 3rem;
		}

		th,
		td {
			padding: 0.75rem;
		}
	}
</style>
