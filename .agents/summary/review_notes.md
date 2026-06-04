# Review Notes

Findings from the consistency/completeness review of `.agents/summary/`.

## Consistency

- **Pipeline framing is consistent** across `codebase_info.md`, `architecture.md`, and `workflows.md` (3-component pipeline communicating via `genFeatureDef.md`).
- **Construct vocabulary** is used consistently in `data_models.md`, `interfaces.md`, and `components.md` and matches `docs/genFeatureDef.md`'s table of contents.
- **Java-abstractor plan** trimmed to 11 remaining steps; docs aligned with `performAbstraction`, shadow→`anyDesc`, partial enums, and `test1006` fixture mismatch noted.
- **Runner stub** is noted in `components.md`, `interfaces.md`, and the index. No drift.
- **CLI flag descriptions** for goAbstractor were derived from `goAbstractor/main.go` (the `argObject` struct). For javaAbstractor, only the *runtime fields* on `Config` consumed in `App.java` are described — the literal flag names are owned by `Config` (not read in this analysis). If you want CLI-flag accuracy verified, ask for a follow-up that reads `Config.java`.

## Completeness Gaps

These are areas where the documentation is intentionally lighter than the source, and where the underlying file is a better reference:

1. **`docs/genFeatureDef.md` schema details.** The summary lists construct kinds and groupings. Per-field schemas (e.g. exact field names in `Method`, `StructDesc`, etc.) are not duplicated. Consult the schema doc directly.
2. **`Config.java` CLI flag names.** Not enumerated; only the fields used in `App.java` are mentioned.
3. **`techDebtMetrics` deep internals.** The .NET projects are described at the file/project level. Specific algorithms in `DesignRecovery/` and `TechDebt/` are not summarized in detail because the `Runner` is stubbed and these libraries are not yet end-to-end exercised. Once the Java abstractor lands and the runner is wired up, this section deserves a refresh.
4. **goAbstractor resolver passes.** Each resolver sub-package (`dce`, `genInterfaces`, `inheritance`, `instantiations`, `references`) is named but not described in detail. For deep work in this area, read the Go source directly.
5. **Test fixture contents.** `testData/go/test*` and `testData/java/test*` are listed but individual fixtures are not described. The plan file's Step descriptions are the better reference for which fixture exercises what.
6. **CI matrix details for techDebtMetrics.** Only a partial slice of `.github/workflows/ci.yaml` was inspected. The Go and Java jobs are documented; the .NET job's matrix specifics are summarized but not field-by-field.
7. **TDD database schema.** Mentioned with a pointer to `.agents/planning/.../research/td-dataset.md`. Not duplicated here.

## Recommendations

- **Refresh after Step 3 lands** (enum completion). Update `data_models.md` and `workflows.md` to mention `Value` constructs being emitted for enum constants and any structural changes to `ObjectDecl` for `CtEnum`.
- **Refresh when the Resolver pipeline is extracted** (plan Step 8). Update `architecture.md` / `workflows.md` with resolver sub-steps when `Resolver.java` exists.
- **Re-run the SOP after the Runner is implemented**, replacing the "stub" notes in `components.md`, `interfaces.md`, and `index.md` with the real CLI surface.
- **Consider a `glossary.md`** if more researchers join — terms like `participation`, `WMC`, `TCC`, `ATFD`, `Cmp`, `pinning`, and the construct names benefit from a single short reference.
- **`AGENTS.md` will be reworked** by the researcher after this run (per the user message). When the new structure stabilizes, this index and `workflows.md` should be updated to reflect any new agent rules (e.g. file-modification scopes mentioned in chat: free-edit `.agents/`, `.cursor/`, `AGENTS.md`; ask-first elsewhere; never `git add`/`commit`/`push`).

## Risks of Staleness

The documents most likely to drift first:
- **`workflows.md`** — Java-abstractor flow will change when the resolver phase is extracted.
- **`components.md`** — refresh when `ObjectDecl.nest`, package `Value`s, or JDK stub cache land.
- **`interfaces.md`** — Runner CLI will replace the stub note.

Auto-generated metrics (file counts, line counts) are intentionally avoided so they don't go stale.
