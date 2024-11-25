// aldff but as a latch
module dllatch #(
    parameter N = 1
) (
    input  logic write, load,
    input  logic [N-1:0] data_in, load_in,
    output logic [N-1:0] data_out,
);
    always begin
        if (load == 1)
            data_out <= load_in;
        else if (write == 1)
            data_out <= data_in;
    end
endmodule