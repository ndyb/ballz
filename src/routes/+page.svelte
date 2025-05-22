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
			const response = await fetch('/api/scores');

			if (!response.ok) {
				throw new Error(`API error: ${response.status}`);
			}

			const scores: Score[] = await response.json();

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

<div class="container mx-auto px-4 py-8">
	<!-- Page title -->
	<h1 class="text-4xl font-bold text-center my-8 text-blue-600">Ballz Leaderboard</h1>

	{#if error}
		<div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-8">
			<p>Error loading scores: {error}</p>
		</div>
	{/if}

	{#if loading}
		<div class="flex justify-center my-12">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else}
		<!-- Highest scores section -->
		<div class="flex flex-col md:flex-row justify-center items-center gap-8 mb-12">
			<!-- All-time high score -->
			<div class="bg-blue-100 rounded-lg p-6 shadow-lg text-center">
				<h2 class="text-xl font-semibold text-blue-800 mb-2">All-Time High</h2>
				{#if allTimeHigh}
					<p class="text-6xl font-bold text-blue-600">{allTimeHigh.score.toLocaleString()}</p>
					<p class="text-gray-600 mt-2">Player: {allTimeHigh.player}</p>
					<p class="text-gray-500 text-sm">{formatDate(allTimeHigh.timestamp)}</p>
				{:else}
					<p class="text-2xl text-gray-500">No data available</p>
				{/if}
			</div>

			<!-- Daily high score -->
			<div class="bg-green-100 rounded-lg p-6 shadow-lg text-center">
				<h2 class="text-xl font-semibold text-green-800 mb-2">Today's High</h2>
				{#if dailyHigh}
					<p class="text-5xl font-bold text-green-600">{dailyHigh.score.toLocaleString()}</p>
					<p class="text-gray-600 mt-2">Player: {dailyHigh.player}</p>
					<p class="text-gray-500 text-sm">{formatDate(dailyHigh.timestamp)}</p>
				{:else}
					<p class="text-2xl text-gray-500">No scores today</p>
				{/if}
			</div>
		</div>

		<!-- Recent scores table -->
		<div class="overflow-x-auto bg-white rounded-lg shadow-md">
			<h2 class="text-2xl font-bold text-center my-4 text-gray-800">Recent Scores</h2>
			<table class="min-w-full">
				<thead class="bg-gray-100">
					<tr>
						<th class="py-3 px-4 text-left text-gray-700 font-semibold">Rank</th>
						<th class="py-3 px-4 text-left text-gray-700 font-semibold">Player</th>
						<th class="py-3 px-4 text-left text-gray-700 font-semibold">Score</th>
						<th class="py-3 px-4 text-left text-gray-700 font-semibold">Date</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-200">
					{#each recentScores as score, i}
						<tr class="hover:bg-gray-50">
							<td class="py-3 px-4 text-gray-800">{i + 1}</td>
							<td class="py-3 px-4 text-gray-800">{score.player}</td>
							<td class="py-3 px-4 text-gray-800 font-semibold">{score.score.toLocaleString()}</td>
							<td class="py-3 px-4 text-gray-600">{formatDate(score.timestamp)}</td>
						</tr>
					{:else}
						<tr>
							<td colspan="4" class="py-4 text-center text-gray-500">No recent scores available</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
