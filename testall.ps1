$fail = 0

go list -m -f '{{.Dir}}' | ForEach-Object {
    $dir = $_.Trim()
    if ($dir -ne "") {
        Write-Host "===> go test -count=1 $dir"
        Push-Location $dir
        try {
            go test -count=1 ./...
        } catch {
            $fail = 1
        } finally {
            Pop-Location
        }
        Write-Host "===> $dir test complete\n"
    }
}

exit $fail

