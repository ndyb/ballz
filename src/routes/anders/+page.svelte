<script lang="ts">
	import { onMount } from 'svelte';

	let scores: number[] = [];
	let currentPage = 1;
	const pageSize = 20;

	let paginatedScores: number[] = [];

	const apiUrl = 'https://your-api.com/scores'; // Replace with your real API

	async function fetchScores() {
		try {
			const response = await fetch(apiUrl);
			const data = await response.json();

			// Sort descending and extract top scores
			scores = data.map((item: { score: number }) => item.score).sort((a, b) => b - a); // Optional: sort high to low

			updatePage();
		} catch (error) {
			console.error('Error fetching scores:', error);
		}
	}

	function updatePage() {
		const start = (currentPage - 1) * pageSize;
		const end = start + pageSize;
		paginatedScores = scores.slice(start, end);
	}

	function nextPage() {
		if (currentPage * pageSize < scores.length) {
			currentPage++;
			updatePage();
		}
	}

	function prevPage() {
		if (currentPage > 1) {
			currentPage--;
			updatePage();
		}
	}

	onMount(() => {
		fetchScores();
	});
</script>

<h2>Highscores</h2>

<table>
	<thead>
		<tr>
			<th>Score</th>
		</tr>
	</thead>
	<tbody>
		{#each paginatedScores as score}
			<tr>
				<td>{score}</td>
			</tr>
		{/each}
	</tbody>
</table>

<div class="pagination">
	<button on:click={prevPage} disabled={currentPage === 1}>Previous</button>
	<span>Page {currentPage}</span>
	<button on:click={nextPage} disabled={currentPage * pageSize >= scores.length}>Next</button>
</div>

<style>
	table {
		width: 100%;
		border-collapse: collapse;
	}

	th,
	td {
		padding: 0.75rem;
		border: 1px solid #ddd;
		text-align: left;
	}

	.pagination {
		margin-top: 1rem;
		display: flex;
		gap: 1rem;
	}

	button:disabled {
		opacity: 0.5;
	}
</style>
