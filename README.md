# <img src="https://github.com/yarox24/EvtxHussar/blob/447cd68ab8f3a4e5bd9d0197461d81cc162b8202/icon/icons8-forensics-96.png" alt="Icon" width="40"/> EvtxHussar



### Input data
- .evtx files coming from various or single host

### Output data
- Subset of events based on event ID's defined in maps (e.g. System 104 - The log file was cleared.)
- Events useful for forensics
- One of the following output formats: CSV, JSON, JSONL, Excel
- Default output format Excel
- Files with the same computer name are merged

#### Example output
###### Subset of columns only (Click for fullscreen preview)
![image](https://user-images.githubusercontent.com/18016218/164982801-4fdc2786-0bfb-439a-8679-1ab35537e4c0.png)

###### Output directory structure (Click for fullscreen preview)
![image](https://user-images.githubusercontent.com/18016218/164982810-7b706507-8bcb-42be-aba5-5bfec0846154.png)

### Interesting features
- Reconstruction of PowerShell Scriptblocks
- Powershell -enc <base64 string> is automatically decoded
- Scheduled Tasks XML parsing
- Merge events from different sources (e.g. Microsoft-Windows-PowerShellOperational_General and Windows PowerShell) to single output file
- Deduplication of events (so you can provide logs from backup, VSS, archive)
- Supported events can be easily added by adding .yaml files to maps/ directory
- Parameters resolution (e.g. %%1936 changed to TokenElevationTypeDefault (1))
- Fields resolution (e.g. servicestarttype = 2 is replaced with "Auto start")

### Which events are supported?
Please look into [maps/](https://github.com/yarox24/EvtxHussar/tree/main/maps "L1 maps") (which contains Layer 1 maps)

### Quick usage

**Parse events (C:\\evtx_compromised_machine\\\*.evtx) from single host to default Excel format**
```cmd
EvtxHussar.exe -o C:\evtxhussar_results C:\evtx_compromised_machine
```

**Parse events (C:\\evtx_many_machines\\\*\\\*.evtx) from many machines recursively saving them with JSONL format**
```cmd
EvtxHussar.exe -f jsonl -r -o C:\evtxhussar_results C:\evtx_many_machines
```

**Parse only 2 files (Security.evtx and System.evtx) and save them with CSV format**
```cmd
EvtxHussar.exe -f csv -o C:\evtxhussar_results C:\evtx_compromised_machine\Security.evtx C:\evtx_compromised_machine\System.evtx
```

**Parse events with 100 workers (1 worker = 1 Evtx file handled) Default: 30**
```cmd
EvtxHussar.exe -w 100 -r -o C:\evtxhussar_results C:\evtx_many_machines
```

### Help
```cmd
Usage: EvtxHussar1.0_amd64.exe [--recursive] [--output_dir OUTPUT_DIR] [--format FORMAT] [--workers WORKERS] [--maps MAPS] [--debug] [INPUT_EVTX_PATHS [INPUT_EVTX_PATHS ...]]

Positional arguments:
  INPUT_EVTX_PATHS       Path(s) to .evtx files or directories containing these files (can be mixed)

Options:
  --recursive, -r        Recursive traversal for any input directories. [default: false]
  --output_dir OUTPUT_DIR, -o OUTPUT_DIR
                         Reports will be saved in this directory (if doesn't exists it will be created)
  --format FORMAT, -f FORMAT
                         Output data in one of the formats: Csv,JSON,JSONL,Excel [default: Excel]
  --workers WORKERS, -w WORKERS
                         Max concurrent workers (.evtx opened) [default: 30]
  --maps MAPS, -m MAPS   Custom directory with maps/ (Default: program directory)
  --debug, -d            Be more verbose [default: false]
  --help, -h             display this help and exit
  --version              display version and exit
```
  
  ### Winged Hussars
  [![Winged Hussars](https://yt-embed.herokuapp.com/embed?v=rcYhYO02f98)](https://www.youtube.com/watch?v=rcYhYO02f98 "Winged Hussars")

  
